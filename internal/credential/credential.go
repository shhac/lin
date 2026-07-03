package credential

import (
	"os"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/credential/keychain"
)

const keychainPlaceholder = "__KEYCHAIN__"

// getKeychain is the keychain lookup, indirected through a package var so a test
// can stub the placeholder-hit branch without touching a real keychain.
var getKeychain = keychain.Get

// workspaceKey decodes a stored workspace record into its API key and source
// ("keychain" or "config"). A keychain placeholder that misses (secret gone or
// unreadable) returns an empty key — never a fallback to another source, so a
// caller naming this alias gets that identity or nothing.
func workspaceKey(alias string, ws config.Workspace) (key, source string) {
	if ws.APIKey == keychainPlaceholder {
		if k, err := getKeychain(alias); err == nil && k != "" {
			return k, "keychain"
		}
		return "", ""
	}
	if ws.APIKey != "" {
		return ws.APIKey, "config"
	}
	return "", ""
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

	key, source := workspaceKey(ws, workspace)
	if key == "" {
		return nil
	}
	return &resolved{key: key, source: source}
}
