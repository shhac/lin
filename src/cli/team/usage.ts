import type { Command } from "commander";

const USAGE_TEXT = `lin team — Manage Linear teams and discover workflow configuration

SUBCOMMANDS:
  team list                  List all teams (id, name, key)
  team get <id>              Team details, members, and estimate config
  team states <team>         List workflow states (valid status values for issue updates)

ARGUMENTS:
  <id>      Team UUID or key (e.g., "ENG")
  <team>    Team key or name (e.g., "ENG" or "Engineering")

OPTIONS (list):
  --limit <n>               Limit results
  --cursor <token>          Pagination cursor for next page

OUTPUT FIELDS:
  list    → id, name, key
  get     → id, name, key, description, estimates (type, allowZero, extended,
            default, validValues, display), members (id, name, email)
  states  → id, name, type, color, position (sorted by position)

ESTIMATE TYPES: notUsed | exponential | fibonacci | linear | tShirt
  Use "team get <id>" to see valid estimate values for a specific team.

WORKFLOW STATE TYPES: triage | backlog | unstarted | started | completed | canceled
  Use "team states <team>" to discover valid status names for issue updates.
  State names are team-specific (e.g., "In Progress", "Todo", "Done").
`;

export function registerUsage(team: Command): void {
  team
    .command("usage")
    .description("Print detailed team command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
