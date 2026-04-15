import type { Command } from "commander";

const USAGE_TEXT = `lin user — Look up Linear workspace users

SUBCOMMANDS:
  user search <text>         Search users by name, email, or display name
  user list [--team <team>]  List users (optionally filtered by team)
  user me                    Current authenticated user + organization

OPTIONS (search):
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OPTIONS (list):
  --team <team>             Filter by team key or name (e.g., "ENG")
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OUTPUT FIELDS:
  search → id, name, email, displayName
  list   → id, name, email, displayName
  me     → id, name, email, displayName, organization (id, name)

NOTES:
  User IDs are UUIDs. Use "user me" to get the current user's ID for
  operations like --assignee filtering.
  When --team is specified, returns only members of that team.
  Search matches against name, email, and displayName (case-insensitive).
`;

export function registerUsage(user: Command): void {
  user
    .command("usage")
    .description("Print detailed user command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
