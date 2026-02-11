import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";

export function registerSearch(project: Command): void {
  project
    .command("search")
    .description("Full-text search for projects")
    .argument("<text>", "Search text")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (text: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const results = await client.searchProjects(text, {
          first: resolvePageSize(opts),
          after: opts.cursor,
        });
        printPaginated(
          results.nodes.map((p) => ({
            id: p.id,
            slugId: p.slugId,
            url: p.url,
            name: p.name,
            status: p.state,
            progress: p.progress,
          })),
          results.pageInfo,
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "Search failed");
      }
    });
}
