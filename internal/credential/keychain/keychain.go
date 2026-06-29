// Package keychain wraps the shared creds.Keychain backend, keyed by lin's
// reverse-domain service. The exported funcs keep their historical
// (string, error) / error signatures so callers in internal/credential and
// internal/config don't change.
package keychain

import (
	"errors"
	"runtime"

	"github.com/shhac/lib-agent-cli/creds"
)

const service = "app.paulie.lin"

var kc = creds.NewKeychain(service)

// MCPKeychainService is the Keychain service for the MCP server's local-OAuth
// secrets — the CLI's service plus a ".mcp" namespace, separate from the API creds.
func MCPKeychainService() string { return service + ".mcp" }

// errNotFound preserves the legacy contract where a missing entry is reported
// as an error (callers check `err == nil && key != ""`).
var errNotFound = errors.New("keychain: entry not found")

// Get retrieves a keychain entry. macOS only; returns "" on other platforms.
func Get(account string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", nil
	}
	v, ok := kc.Get(account)
	if !ok {
		return "", errNotFound
	}
	return v, nil
}

// Store saves a keychain entry. macOS only; no-op on other platforms.
func Store(account, password string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	return kc.Set(account, password)
}

// Delete removes a keychain entry. macOS only; no-op on other platforms.
func Delete(account string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	return kc.Delete(account)
}

// DeleteAll removes every entry for the service, including orphans not tracked
// in config. macOS only; no-op elsewhere.
func DeleteAll() {
	if runtime.GOOS != "darwin" {
		return
	}
	_ = kc.DeleteAll()
}
