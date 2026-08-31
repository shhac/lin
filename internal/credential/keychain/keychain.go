// Package keychain wraps the shared creds.Keychain backend, keyed by lin's
// reverse-domain service. The exported funcs keep their historical
// (string, error) / error signatures so callers in internal/credential and
// internal/config don't change.
//
// Platform support and the LIN_NO_KEYCHAIN / LIB_AGENT_NO_KEYCHAIN opt-out are
// the shared backend's job, not this wrapper's: it has a real store on macOS,
// Linux and Windows, and reports itself unavailable elsewhere. Do not reinstate
// a runtime.GOOS gate here — short-circuiting above the backend skips those
// stores and, worse, turns "no keychain" into a silent successful no-op, which
// makes StoreLogin swap a live secret for the placeholder and lose it.
package keychain

import (
	"errors"

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

// Get retrieves a keychain entry, reporting a miss — including a keychain that
// is unavailable or opted out — as errNotFound.
func Get(account string) (string, error) {
	v, ok := kc.Get(account)
	if !ok {
		return "", errNotFound
	}
	return v, nil
}

// Store saves a keychain entry. An unavailable keychain returns an error rather
// than reporting a phantom success, so callers fall back to the file store.
func Store(account, password string) error {
	return kc.Set(account, password)
}

// Delete removes a keychain entry.
func Delete(account string) error {
	return kc.Delete(account)
}

// DeleteAll removes every entry for the service, including orphans not tracked
// in config.
func DeleteAll() {
	_ = kc.DeleteAll()
}
