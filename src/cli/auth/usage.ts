import type { Command } from "commander";

const USAGE_TEXT = `lin auth — Authentication and workspace management

SUBCOMMANDS:
  auth login <api-key> [--alias <name>]   Store API key and detect workspace
  auth status                             Show current auth state
  auth logout [--all]                     Remove stored credentials
  auth workspace list                     List all stored workspaces
  auth workspace switch <alias>           Set default workspace
  auth workspace remove <alias>           Remove a stored workspace

AUTH SOURCES (checked in order):
  1. LINEAR_API_KEY environment variable (takes precedence)
  2. Stored config (~/.config/lin/config.json)

LOGIN:
  Validates the API key against the Linear API, auto-detects org name/urlKey.
  --alias sets a custom workspace name (default: org urlKey).
  Multiple workspaces supported — each login adds a workspace profile.

LOGOUT:
  Default: removes only the active workspace, switches default to next available.
  --all: clears all stored workspaces and credentials.

STATUS:
  Returns { authenticated, source, user, organization, activeWorkspace, otherWorkspaces }.
  source is "environment" if LINEAR_API_KEY is set, "config" otherwise.

WORKSPACE:
  list — shows all stored workspaces with alias, name, urlKey, and which is default.
  switch <alias> — sets the active workspace. Alias must match a stored workspace.
  remove <alias> — deletes a workspace profile. Warns if removing the default.

OUTPUT: JSON to stdout. Errors: { "error": "..." } to stderr.
`;

export function registerUsage(auth: Command): void {
  auth
    .command("usage")
    .description("Print detailed auth command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
