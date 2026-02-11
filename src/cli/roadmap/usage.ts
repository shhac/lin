import type { Command } from "commander";

const USAGE_TEXT = `lin roadmap — Browse Linear roadmaps (read-only)

LIST:
  roadmap list                            List all roadmaps
    [--limit] [--cursor]
    Returns per item: id, slugId, url, name, description, owner

GET:
  roadmap get overview <id>    Roadmap summary: id, slugId, url, name, description,
                               owner, creator, createdAt
  roadmap get projects <id>    Projects linked to a roadmap
    [--limit] [--cursor]
    Returns per item: id, slugId, url, name, status, progress,
    lead, startDate, targetDate

IDS: <id> accepts UUID, slug ID, or roadmap name.
PAGINATION: --limit <n> --cursor <token> on list and get projects.

NOTE: Roadmaps are read-only. Use "project" commands to modify projects in a roadmap.
`;

export function registerUsage(roadmap: Command): void {
  roadmap
    .command("usage")
    .description("Print detailed roadmap command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
