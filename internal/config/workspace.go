package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shhac/lin/internal/credential/keychain"
)

const keychainPlaceholder = "__KEYCHAIN__"

// StoreLogin stores a workspace login. Attempts keychain storage first;
// falls back to plaintext config. Clears legacy api_key to avoid stale fallback.
//
// The keychain write happens before the locked config update: by the time the
// config write starts, the secret is already wherever it is going to live, so
// the only step that must be serialized against other lin processes is the
// config entry that references it.
func StoreLogin(alias string, ws Workspace) error {
	if err := keychain.Store(alias, ws.APIKey); err == nil {
		ws.APIKey = keychainPlaceholder
	}

	return update(func(cfg *Config) error {
		cfg.LegacyAPIKey = ""
		if cfg.Workspaces == nil {
			cfg.Workspaces = make(map[string]Workspace)
		}
		cfg.Workspaces[alias] = ws
		if cfg.DefaultWorkspace == "" {
			cfg.DefaultWorkspace = alias
		}
		return nil
	})
}

// SetDefaultWorkspace sets the active workspace alias.
func SetDefaultWorkspace(alias string) error {
	return update(func(cfg *Config) error {
		if _, ok := cfg.Workspaces[alias]; !ok {
			return workspaceNotFoundError(alias, cfg)
		}
		cfg.DefaultWorkspace = alias
		return nil
	})
}

// RemoveWorkspace removes a stored workspace and its keychain entry.
func RemoveWorkspace(alias string) error {
	return update(func(cfg *Config) error {
		if _, ok := cfg.Workspaces[alias]; !ok {
			return workspaceNotFoundError(alias, cfg)
		}
		_ = keychain.Delete(alias)
		delete(cfg.Workspaces, alias)
		if cfg.DefaultWorkspace == alias {
			remaining := workspaceKeys(cfg)
			if len(remaining) > 0 {
				cfg.DefaultWorkspace = remaining[0]
			} else {
				cfg.DefaultWorkspace = ""
			}
		}
		if len(cfg.Workspaces) == 0 {
			cfg.Workspaces = nil
		}
		return nil
	})
}

// ClearApiKey removes the legacy top-level api_key.
func ClearApiKey() error {
	return update(func(cfg *Config) error {
		cfg.LegacyAPIKey = ""
		return nil
	})
}

// ClearAll removes all workspaces and credentials.
func ClearAll() error {
	keychain.DeleteAll()
	return update(func(cfg *Config) error {
		*cfg = Config{Workspaces: make(map[string]Workspace)}
		return nil
	})
}

// GetDefaultWorkspace returns the default workspace alias.
func GetDefaultWorkspace() string {
	return Read().DefaultWorkspace
}

// GetWorkspaces returns all stored workspaces.
func GetWorkspaces() map[string]Workspace {
	cfg := Read()
	if cfg.Workspaces == nil {
		return make(map[string]Workspace)
	}
	return cfg.Workspaces
}

func workspaceKeys(cfg *Config) []string {
	keys := make([]string, 0, len(cfg.Workspaces))
	for k := range cfg.Workspaces {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func workspaceNotFoundError(alias string, cfg *Config) error {
	keys := workspaceKeys(cfg)
	valid := "(none)"
	if len(keys) > 0 {
		valid = strings.Join(keys, ", ")
	}
	return fmt.Errorf("unknown workspace: %s, valid: %s", alias, valid)
}
