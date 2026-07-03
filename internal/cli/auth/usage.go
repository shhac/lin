package auth

const usageText = `lin auth — Authentication and workspace management

SUBCOMMANDS:
  auth login <api-key> [--alias <name>]   Store API key and detect workspace
  auth status                             Show current auth state
  auth logout [--all]                     Remove stored credentials
  auth workspace list                     List all stored workspaces
  auth workspace switch <alias>           Set default workspace
  auth workspace remove <alias>           Remove a stored workspace

AUTH SOURCES (checked in order):
  1. LINEAR_API_KEY environment variable (takes precedence)
  2. macOS Keychain (service: app.paulie.lin)
  3. Stored config (~/.config/lin/config.json)

WORKSPACE SELECTION:
  --workspace <alias> (global flag) acts as a specific stored workspace for one
  invocation, overriding the default. It resolves strictly by alias against the
  stored workspaces (never env or the default), so an unknown alias errors.

  LIN_REQUIRE_IDENTITY (env, fail-closed): when set, any command that would
  touch Linear WITHOUT an explicit --workspace errors before any fallback
  (default workspace, legacy api_key, or LINEAR_API_KEY) can serve it. The MCP
  server sets this for every named-principal tool call.

LOGIN:
  Validates the API key against the Linear API, auto-detects org name/urlKey.
  --alias sets a custom workspace name (default: org urlKey).
  Multiple workspaces supported — each login adds a workspace profile.

LOGOUT:
  Default: removes only the active workspace, switches default to next available.
  --all: clears all stored workspaces and credentials.

STATUS:
  Returns { authenticated, source, user, organization, activeWorkspace, otherWorkspaces }.
  source is "environment", "keychain", or "config".

WORKSPACE:
  list — shows all stored workspaces with alias, name, urlKey, and which is default.
  switch <alias> — sets the active workspace. Alias must match a stored workspace.
  remove <alias> — deletes a workspace profile. Warns if removing the default.

OUTPUT: workspace list follows list output defaults (JSONL unless --format is set). Other auth commands return JSON objects. Errors: { "error": "..." } to stderr.`
