import type { Command } from "commander";

const USAGE_TEXT = `lin — Linear CLI

COMMANDS:
  auth login <api-key>                    Store API key
  auth status                             Show auth state + workspace

  project search <text>                   Search projects
  project list [--team] [--status]        List projects
  project get overview <id>               Project summary + milestones
  project get details <id>               Full project content (markdown)
  project get issues <id>                Issues in a project
  project update title <id> <value>       Update title
  project update status <id> <value>      Update status
  project update description <id> <value> Update description
  project update lead <id> <user-id>      Update lead

  issue search <text>                     Full-text search
  issue list [filters]                    List issues
  issue get overview <id>                 Issue details
  issue get comments <id>                 List comments
  issue new <title> --team <team>         Create issue
  issue update title <id> <value>         Update title
  issue update status <id> <value>        Update status
  issue update assignee <id> <user-id>    Update assignee
  issue update priority <id> <priority>   Update priority
  issue update project <id> <project-id>  Move to project
  issue update labels <id> <l1,l2,...>    Set labels
  issue update description <id> <value>   Update description
  issue comment new <issue-id> <body>     Add comment

  team list                               List teams
  team get <id>                           Team details + members

  user list [--team]                      List users
  user me                                 Current user

  label list [--team]                     List labels

  cycle list --team <team>                List cycles
  cycle get <id>                          Cycle details

ID FORMAT:
  Issue keys (ENG-123), UUIDs, or URL fragments accepted.

FILTERS:
  --team, --status, --assignee, --priority, --label, --cycle, --limit, --project

OUTPUT:
  All output is JSON to stdout. Errors go to stderr.
  Lists return arrays. Single items return objects.

PRIORITY VALUES:
  none | urgent | high | medium | low

AUTH:
  Set LINEAR_API_KEY env var, or run: lin auth login <key>
`;

export function registerUsageCommand({ program }: { program: Command }): void {
  program
    .command("usage")
    .description("Print concise documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
