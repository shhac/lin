import type { Command } from "commander";

const USAGE_TEXT = `lin label — List Linear issue labels

SUBCOMMANDS:
  label list [--team <team>]  List labels (optionally filtered by team)

OPTIONS:
  --team <team>             Filter by team key or name (e.g., "ENG")
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OUTPUT FIELDS:
  list → id, name, color

NOTES:
  Without --team, returns all workspace labels.
  With --team, returns only labels scoped to that team.
  Label names are used as values for --label filters and "issue update labels".
`;

export function registerUsage(label: Command): void {
  label
    .command("usage")
    .description("Print detailed label command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
