package auth

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

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

OUTPUT: JSON to stdout. Errors: { "error": "..." } to stderr.`

func registerUsage(auth *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Print detailed auth command documentation (LLM-optimized)",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(strings.TrimSpace(usageText))
		},
	}
	auth.AddCommand(cmd)
}
