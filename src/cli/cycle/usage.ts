import type { Command } from "commander";

const USAGE_TEXT = `lin cycle — List and inspect Linear cycles (sprints)

SUBCOMMANDS:
  cycle list --team <team>   List cycles for a team (--team is required)
  cycle get <id>             Cycle details with all issues

ARGUMENTS:
  <id>      Cycle UUID

OPTIONS (list):
  --team <team>             Team key or name (required)
  --current                 Show only the active cycle
  --next                    Show only the upcoming cycle
  --previous                Show only the most recently completed cycle
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OUTPUT FIELDS:
  list → id, number, name, startsAt, endsAt
  get  → id, number, name, startsAt, endsAt,
         issues (id, identifier, title, status, assignee, priority, priorityLabel)

NOTES:
  --current, --next, and --previous are mutually exclusive convenience filters.
  --current returns the team's active cycle (may be empty array if none active).
  --next returns the nearest future cycle (startsAt > now).
  --previous returns the most recently ended cycle (endsAt < now).
  Cycle IDs are UUIDs. Use "cycle list --team ENG" to find cycle IDs.
`;

export function registerUsage(cycle: Command): void {
  cycle
    .command("usage")
    .description("Print detailed cycle command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
