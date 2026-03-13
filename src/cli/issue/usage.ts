import type { Command } from "commander";

const USAGE_TEXT = `lin issue — Issue operations (search, list, create, update, comment, relate, archive, attach, history)

SEARCH & LIST:
  issue search <text> [--project] [--team] [--assignee] [--status] [--priority]
  issue list [--project] [--team] [--assignee] [--status] [--priority] [--label] [--cycle]

GET:
  issue get <id>             Full details + labels, attachments, branchName, commentCount

CREATE:
  issue new <title> --team <key|name|UUID>    --team required
    [--project <name|slug|UUID>] [--assignee] [--priority] [--status]
    [--labels <ids>] [--description <md>] [--cycle] [--parent <id>] [--estimate <n>]

UPDATE (each is a subcommand):
  issue update title|status|assignee|priority|project|labels|description|estimate <id> <value>
  issue update due-date <id> <YYYY-MM-DD>
  issue update cycle <id> <cycle-id>  |  parent <id> <parent-id>
  status: team-scoped (use "team states <team>"). labels: comma-separated IDs.
  estimate: validated against team scale. assignee: name/email/ID.

COMMENTS (--file repeatable, --parent 1 level):
  issue comment list <id> [--limit] [--cursor]
  issue comment new <id> <body> [--parent <cid>] [--file <path>]
  issue comment edit <cid> <body> [--file <path>]
  issue comment get <cid>  |  replies <cid>

RELATIONS:
  issue relation list <id>      Both directions (blocks, blocked_by, duplicate, related)
  issue relation add <id> --type <t> --related <id>  Types: blocks|duplicate|related
  issue relation remove <id>

HISTORY:
  issue history <id> [--limit] [--cursor]     Activity log (status, assignee, priority, labels, etc.)

LIFECYCLE:  issue archive|unarchive|delete <id>

ATTACHMENTS:
  issue attachment list|remove <id>
  issue attachment add <id> --url <u> --title <t> [--subtitle <s>]

IDS: Issue keys (ENG-123) or UUIDs. PRIORITY: none|urgent|high|medium|low
PAGINATION: --limit <n> --cursor <token> on search and list.
`;

export function registerUsage(issue: Command): void {
  issue
    .command("usage")
    .description("Print detailed issue command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
