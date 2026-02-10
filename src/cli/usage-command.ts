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
  project update status <id> <value>      Update status
  project update description <id> <value> Update description
  project update lead <id> <user-id>      Update lead
  project new <name> --team <teams>       Create project (comma-separated team keys/IDs)

  roadmap list                            List roadmaps
  roadmap get overview <id>               Roadmap summary + owner
  roadmap get projects <id>               Projects in a roadmap

  issue search <text>                     Full-text search
  issue list [filters]                    List issues (status, assignee, team, branchName)
  issue get overview <id>                 Full details + commentCount, attachments, branchName
  issue get comments <id>                 List comments with authors
  issue new <title> --team <team>         Create issue (team key or ID)
  issue update title <id> <value>         Update title
  issue update status <id> <value>        Update status (team-specific workflow states)
  issue update assignee <id> <user>       Update assignee (name, email, or ID)
  issue update priority <id> <priority>   Update priority
  issue update project <id> <project>     Move to project (ID, slug, or name)
  issue update labels <id> <l1,l2,...>    Set labels
  issue update estimate <id> <value>      Update estimate (validated against team scale)
  issue update description <id> <value>   Update description
  issue comment new <issue-id> <body>     Add comment
  issue comment get <comment-id>          Get a specific comment
  issue comment edit <comment-id> <body>  Edit a comment
  issue relation list <issue-id>          List relations (both directions)
  issue relation add <id> --type <t> --related <id>  Add relation (blocks|duplicate|related)
  issue relation remove <relation-id>     Remove relation
  issue archive <id>                      Archive issue
  issue unarchive <id>                    Unarchive issue
  issue delete <id>                       Delete issue (trash)
  issue attachment list <issue-id>        List attachments
  issue attachment add <id> --url <u> --title <t>    Add URL attachment
  issue attachment remove <attachment-id> Remove attachment

  team list                               List teams (id, name, key)
  team get <id>                           Team details + members + estimate config
  team states <team>                      List workflow states (valid status values)

  user list [--team]                      List users
  user me                                 Current user

  label list [--team]                     List labels
  cycle list --team <team>                List cycles
  cycle get <id>                          Cycle details

IDS: Issue keys (ENG-123), UUIDs, or URL slugs accepted.
     --team accepts team key (ENG) or name.
     --project/project/roadmap <id> accept UUID, slug ID, or name.

FILTERS (list/search): --team --status --assignee --priority --label --cycle --project --limit

PAGINATION: --limit <n> --cursor <token>
  Output: { "items": [...], "pagination": { "hasMore": true, "nextCursor": "..." } }

OUTPUT: JSON to stdout. Errors: { "error": "..." } to stderr with valid values.

PRIORITY: none | urgent | high | medium | low
PROJECT STATUS: backlog | planned | started | paused | completed | canceled
ESTIMATES: Team-specific scales (fibonacci, linear, exponential, tShirt).
  Use "team get <id>" to see valid estimate values for a team.

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
