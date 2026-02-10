import type { Command } from "commander";

const USAGE_TEXT = `lin — Linear CLI (JSON output, LLM-friendly)

COMMANDS:
  auth login <api-key> [--alias <name>]   Store API key (auto-detects workspace)
  auth logout [--all]                     Remove workspace (--all: clear all)
  auth status                             Show auth state
  auth workspace list|switch|remove       Manage workspaces

  project search <text>                   Search projects
  project list [--team] [--status]        List projects
  project get overview <id>               Project summary + milestones + content
  project get issues <id>                 Issues in a project
  project update <field> <id> <value>     Update (title|status|description|lead)
  project new <name> --team <t> [options]  Create project

  roadmap list                            List roadmaps
  roadmap get overview <id>               Roadmap summary + owner
  roadmap get projects <id>               Projects in a roadmap

  document search <text>                  Full-text search
  document list [--project] [--creator]   List documents
  document get <id>                       Full details + content (markdown)
  document new <title>                    Create document
  document update title|content|project <id> <value>  Update field

  issue search <text>                     Full-text search
  issue list [filters]                    List issues (status, assignee, team, branchName)
  issue get overview <id>                 Full details + commentCount, attachments, branchName
  issue get comments <id>                 List comments with authors
  issue new <title> --team <t> [options]  Create issue
  issue update <field> <id> <value>       Update (title|status|assignee|priority|project|labels|estimate|description)
  issue comment new <issue-id> <body>     Add comment
  issue comment get <comment-id>          Get a specific comment
  issue comment edit <comment-id> <body>  Edit a comment
  issue relation list <issue-id>          List relations
  issue relation add <id> --type <t> --related <id>  Add (blocks|duplicate|related)
  issue relation remove <relation-id>     Remove relation
  issue archive|unarchive|delete <id>     Archive, restore, or trash
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

IDS: Issue keys (ENG-123), UUIDs, or slugs. --team accepts key (ENG) or name.
     Project/roadmap/document <id> accept UUID, slug ID, or name.

FILTERS (issue list/search): --team --status --assignee --priority --label --cycle --project --limit
FILTERS (document list): --project --creator --limit

PAGINATION: --limit <n> --cursor <token>
  { "items": [...], "pagination": { "hasMore": true, "nextCursor": "..." } }

OUTPUT: JSON to stdout. Errors: { "error": "..." } to stderr with valid values.

TRUNCATION: description/body/content truncated to ~200 chars + companion *Length field.
  --expand <field,...>  Expand specific    --full  Expand all (e.g. lin --full issue get overview ENG-123)

PRIORITY: none|urgent|high|medium|low
PROJECT STATUS: backlog|planned|started|paused|completed|canceled
ESTIMATES: Team-specific (fibonacci|linear|exponential|tShirt). Use "team get <id>" for values.

AUTH: Set LINEAR_API_KEY env var, or: lin auth login <key>
  Multiple workspaces supported. Switch: lin auth workspace switch <alias>
`;

export function registerUsageCommand({ program }: { program: Command }): void {
  program
    .command("usage")
    .description("Print concise documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
