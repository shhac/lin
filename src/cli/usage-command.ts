import type { Command } from "commander";

const USAGE_TEXT = `lin — Linear CLI (JSON output, LLM-friendly)

COMMANDS:
  auth login <api-key> [--alias <name>]   Store API key (auto-detects workspace)
  auth logout [--all]                     Remove active workspace (--all: clear everything)
  auth status                             Show auth state + active workspace
  auth workspace list                     List stored workspaces
  auth workspace switch <alias>           Set default workspace
  auth workspace remove <alias>           Remove a workspace

  project search <text>                   Search projects
  project list [--team] [--status]        List projects
  project get overview <id>               Project summary + milestones
  project get details <id>                Full project content (markdown)
  project get issues <id>                 Issues in a project
  project update title <id> <value>       Update title
  project update status <id> <value>      Update status (backlog|planned|started|paused|completed|canceled)
  project update description <id> <value> Update description
  project update lead <id> <user-id>      Update lead

  issue search <text>                     Full-text search
  issue list [filters]                    List issues (status, assignee, team, branchName)
  issue get overview <id>                 Full details + commentCount, attachments, branchName
  issue get comments <id>                 List comments with authors
  issue new <title> --team <team>         Create issue (team key or ID)
  issue update title <id> <value>         Update title
  issue update status <id> <value>        Update status (team-specific workflow states)
  issue update assignee <id> <user-id>    Update assignee
  issue update priority <id> <priority>   Update priority
  issue update project <id> <project-id>  Move to project
  issue update labels <id> <l1,l2,...>    Set labels
  issue update description <id> <value>   Update description
  issue comment new <issue-id> <body>     Add comment
  issue comment get <comment-id>          Get a specific comment
  issue comment edit <comment-id> <body>  Edit a comment

  team list                               List teams (id, name, key)
  team get <id>                           Team details + members

  user list [--team]                      List users
  user me                                 Current user

  label list [--team]                     List labels
  cycle list --team <team>                List cycles
  cycle get <id>                          Cycle details

IDS: Issue keys (ENG-123), UUIDs, or URL slugs accepted.
     --team accepts team key (ENG) or name.
     --project accepts project UUID, slug ID, or name.
     Project commands (<id>) accept UUID, slug ID, or name.

FILTERS (list/search): --team --status --assignee --priority --label --cycle --project --limit

PAGINATION: --limit <n> --cursor <token>
  List output: { "items": [...], "pagination": { "hasMore": true, "nextCursor": "..." } }
  When no more pages, pagination key is omitted.

OUTPUT: JSON to stdout. Errors: { "error": "..." } to stderr.
  Error messages include valid values when input is invalid.
  issue get overview includes: commentCount, branchName, attachments (PR links etc.)

PRIORITY: none | urgent | high | medium | low
PROJECT STATUS: backlog | planned | started | paused | completed | canceled

AUTH: Set LINEAR_API_KEY env var, or: lin auth login <key>
  Multiple workspaces: lin auth login <key1>, lin auth login <key2>
  Switch: lin auth workspace switch <alias>
`;

export function registerUsageCommand({ program }: { program: Command }): void {
  program
    .command("usage")
    .description("Print concise documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
