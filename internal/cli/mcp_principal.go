// Per-principal MCP wiring: how an authenticated named principal's pairing
// binding (mcp pair add <name> --bind workspace=<alias>) shapes their tool
// calls — the credential selector plus the fail-closed gate on every
// subprocess. lin deliberately does NOT scope file roots per principal (see the
// WithCredentialEnrollment wiring in root.go): its cache is not identity-
// namespaced, so the lib's default of hiding roots from named principals is the
// safe behavior.
package cli

import (
	"github.com/shhac/lib-agent-mcp/oauth"
)

// bindingKeyWorkspace is the pairing-binding key naming the workspace alias a
// principal acts as — the vocabulary contract with
// `mcp pair add <name> --bind workspace=<alias>`.
const bindingKeyWorkspace = "workspace"

// mcpIdentityBinding translates an authenticated MCP principal into the
// subprocess shape their tool calls run with: `--workspace <alias>` pinning the
// stored credentials their pairing was bound to, plus the fail-closed gate so a
// call that arrives without a selector — a binding-plumbing bug — errors
// instead of falling back to the operator's default workspace. The MCP layer
// only invokes this for named principals; the anonymous operator (stdio, legacy
// shared pairing code) stays unbound.
func mcpIdentityBinding(p oauth.Verified) (argv, env []string) {
	env = []string{"LIN_REQUIRE_IDENTITY=1"}
	alias := p.Binding[bindingKeyWorkspace]
	if alias == "" {
		return nil, env
	}
	return []string{"--workspace", alias}, env
}
