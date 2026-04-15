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
	if key := os.Getenv("LINEAR_API_KEY"); key != "" {
		return key
	}

	cfg := config.Read()
	ws := cfg.DefaultWorkspace
	if ws == "" {
		return cfg.LegacyAPIKey
	}

	workspace, ok := cfg.Workspaces[ws]
	if !ok {
		return cfg.LegacyAPIKey
	}

	if workspace.APIKey == keychainPlaceholder {
		if key, err := keychain.Get(ws); err == nil && key != "" {
			return key
		}
	}

	return workspace.APIKey
}
