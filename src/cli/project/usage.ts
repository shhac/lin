import type { Command } from "commander";

const USAGE_TEXT = `lin project — Project operations (search, list, get, create, update)

SEARCH & LIST:
  project search <text>                   Full-text search for projects
    [--limit] [--cursor]
  project list                            List all projects
    [--team <team>] [--status <status>] [--limit] [--cursor]

GET:
  project get overview <id>    Project summary: id, slugId, url, name, description, content,
                               status, progress, lead, startDate, targetDate, milestones[]
  project get issues <id>      Issues in a project
    [--status] [--assignee] [--priority] [--limit] [--cursor]

CREATE:
  project new <name> --team <teams>       --team required (comma-separated for multi-team)
    [--description <desc>] [--lead <user>] [--start-date <YYYY-MM-DD>]
    [--target-date <YYYY-MM-DD>] [--status <status>] [--content <markdown>]

UPDATE (each is a subcommand):
  project update title <id> <new-title>
  project update status <id> <new-status>
  project update description <id> <description>
  project update lead <id> <user-id>

IDS: <id> accepts UUID, slug ID, or project name.
TEAM: Team key (ENG) or team name.
LEAD: Name, email, or user ID (on create). User ID only (on update).
PROJECT STATUS: backlog|planned|started|paused|completed|canceled
PAGINATION: --limit <n> --cursor <token> on search, list, and get issues.
`;

export function registerUsage(project: Command): void {
  project
    .command("usage")
    .description("Print detailed project command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
