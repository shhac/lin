import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";

export function registerList(roadmap: Command): void {
  roadmap
    .command("list")
    .description("List roadmaps")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: { limit?: string; cursor?: string }) => {
      try {
        const client = getClient();
        const results = await client.roadmaps({
          first: resolvePageSize(opts),
          after: opts.cursor,
        });
        const items = await Promise.all(
          results.nodes.map(async (r) => {
            const owner = await r.owner;
            return {
              id: r.id,
              slugId: r.slugId,
              url: r.url,
              name: r.name,
              description: r.description,
              owner: owner ? owner.name : null,
            };
          }),
        );
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
