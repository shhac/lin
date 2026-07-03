package credential

import (
	"os"
	"sort"
	"strings"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/credential/keychain"
	apierrors "github.com/shhac/lin/internal/errors"
)

const keychainPlaceholder = "__KEYCHAIN__"

// selectedWorkspace pins resolution to one stored workspace alias for the
// current invocation (the --workspace flag). Empty restores the default.
var selectedWorkspace string

// SetSelectedWorkspace records the --workspace selector for this invocation.
// The value is resolved strictly against the workspaces map by ResolveForClient
// (never env or the default), so a caller naming an identity always gets that
// one or an error — never a silent fallback to someone else's.
func SetSelectedWorkspace(alias string) { selectedWorkspace = strings.TrimSpace(alias) }

// ResolveForClient resolves the Linear API key honoring the --workspace
// selector and the LIN_REQUIRE_IDENTITY fail-closed gate. It returns a
// structured error for the gate and for an unknown explicit selector; a plain
// empty string with a nil error still means "no credentials", so callers keep
// their own not-authenticated message.
func ResolveForClient() (string, error) {
	selector := selectedWorkspace

	// Fail-closed identity mode: an MCP runner serving several principals sets
	// LIN_REQUIRE_IDENTITY and passes an explicit --workspace on every
	// invocation. A missing selector then means the caller's identity binding
	// was not applied — refuse before ANY credential source (default workspace,
	// legacy api_key, or LINEAR_API_KEY env) can serve the request as the wrong
	// identity.
	if selector == "" && requireIdentity() {
		return "", apierrors.New(
			"LIN_REQUIRE_IDENTITY is set but no workspace was specified", apierrors.FixableByAgent).
			WithHint("pass --workspace <alias>; falling back to the default workspace is disabled in this environment")
	}

	if selector != "" {
		return resolveSelector(selector)
	}
	return Resolve(), nil
}

// resolveSelector resolves strictly by alias from the workspaces map
// (keychain-first) — never env or the default, since the caller named the
// identity. An unknown alias is an agent-fixable error listing the known ones.
func resolveSelector(alias string) (string, error) {
	cfg := config.Read()
	ws, ok := cfg.Workspaces[alias]
	if !ok {
		return "", unknownWorkspaceError(alias, cfg)
	}
	if ws.APIKey == keychainPlaceholder {
		if key, err := keychain.Get(alias); err == nil && key != "" {
			return key, nil
		}
		return "", nil
	}
	return ws.APIKey, nil
}

func unknownWorkspaceError(alias string, cfg *config.Config) error {
	keys := make([]string, 0, len(cfg.Workspaces))
	for k := range cfg.Workspaces {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	known := "(none)"
	if len(keys) > 0 {
		known = strings.Join(keys, ", ")
	}
	return apierrors.Newf(apierrors.FixableByAgent,
		"unknown workspace %q; known workspaces: %s", alias, known).
		WithHint("pass --workspace <alias> matching a configured workspace, or run 'lin auth login' to add one")
}

// requireIdentity reports whether the fail-closed identity gate is on. Any
// value except empty/0/false counts as set.
func requireIdentity() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LIN_REQUIRE_IDENTITY"))) {
	case "", "0", "false":
		return false
	default:
		return true
	}
}

// Resolve returns the Linear API key, checking in order:
// 1. LINEAR_API_KEY environment variable
// 2. macOS Keychain (if config has placeholder)
// 3. Config file plaintext
// 4. Legacy top-level api_key field
func Resolve() string {
	r := resolve()
	if r == nil {
		return ""
	}
	return r.key
}

// Source returns the source of the resolved API key ("environment", "keychain", or "config").
// Returns empty string if no key is found.
func Source() string {
	r := resolve()
	if r == nil {
		return ""
	}
	return r.source
}

type resolved struct {
	key    string
	source string
}

func resolve() *resolved {
	if key := os.Getenv("LINEAR_API_KEY"); key != "" {
		return &resolved{key: key, source: "environment"}
	}

	cfg := config.Read()
	legacy := func() *resolved {
		if cfg.LegacyAPIKey != "" {
			return &resolved{key: cfg.LegacyAPIKey, source: "config"}
		}
		return nil
	}

	ws := cfg.DefaultWorkspace
	if ws == "" {
		return legacy()
	}
	workspace, ok := cfg.Workspaces[ws]
	if !ok {
		return legacy()
	}

	if workspace.APIKey == keychainPlaceholder {
		if key, err := keychain.Get(ws); err == nil && key != "" {
			return &resolved{key: key, source: "keychain"}
		}
		return nil
	}
	if workspace.APIKey != "" {
		return &resolved{key: workspace.APIKey, source: "config"}
	}

	return nil
}
