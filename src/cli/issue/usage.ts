import type { Command } from "commander";

const USAGE_TEXT = `lin issue — Issue operations (search, list, create, update, comment, relate, archive, attach)

SEARCH & LIST:
  issue search <text> [--project] [--team] [--assignee] [--status] [--priority]
  issue list [--project] [--team] [--assignee] [--status] [--priority] [--label] [--cycle]

GET:
  issue get overview <id>    Full details + labels, attachments, branchName, commentCount
  issue get comments <id>    Comments with author, body, timestamps

CREATE:
  issue new <title> --team <key|name|UUID>    --team required
    [--project <name|slug|UUID>] [--assignee] [--priority] [--status]
    [--labels <ids>] [--description <md>] [--cycle] [--parent <id>] [--estimate <n>]

UPDATE (each is a subcommand):
  issue update title|status|assignee|priority|project|labels|description|estimate <id> <value>
  status: team-scoped name (use "team states <team>")
  assignee: name, email, or user ID
  labels: comma-separated label IDs
  estimate: validated against team scale

COMMENTS (--file repeatable, --parent 1 level):
  issue comment new <id> <body> [--parent <cid>] [--file <path>]
  issue comment edit <cid> <body> [--file <path>]
  issue comment get <cid>      Author, issue, parent, childCount
  issue comment replies <cid> List replies

RELATIONS:
  issue relation list <id>      Both directions (blocks, blocked_by, duplicate, related)
  issue relation add <id> --type <t> --related <id>  Types: blocks|duplicate|related
  issue relation remove <id>

LIFECYCLE:
  issue archive|unarchive|delete <id>      Archive, restore, or trash

ATTACHMENTS:
  issue attachment list <issue-id>
  issue attachment add <id> --url <u> --title <t> [--subtitle <s>]
  issue attachment remove <id>

IDS: Issue keys (ENG-123) or UUIDs. Comment/relation/attachment IDs are UUIDs.
ASSIGNEE FILTER: name, email, user ID, or "me". PRIORITY: none|urgent|high|medium|low
ESTIMATES: Team-specific (fibonacci|linear|exponential|tShirt). Run "team get <id>".
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
