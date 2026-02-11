import type { Command } from "commander";

const USAGE_TEXT = `lin user — Look up Linear workspace users

SUBCOMMANDS:
  user list [--team <team>]  List users (optionally filtered by team)
  user me                    Current authenticated user + organization

OPTIONS (list):
  --team <team>             Filter by team key or name (e.g., "ENG")
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OUTPUT FIELDS:
  list → id, name, email, displayName
  me   → id, name, email, displayName, organization (id, name)

NOTES:
  User IDs are UUIDs. Use "user me" to get the current user's ID for
  operations like --assignee filtering.
  When --team is specified, returns only members of that team.
`;

export function registerUsage(user: Command): void {
  user
    .command("usage")
    .description("Print detailed user command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
