package credential

import (
	"os"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/credential/keychain"
)

const keychainPlaceholder = "__KEYCHAIN__"

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
