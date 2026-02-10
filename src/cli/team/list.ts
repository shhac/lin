import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated } from "../../lib/output.ts";

export function registerList(team: Command): void {
  team
    .command("list")
    .description("List all teams")
    .option("--limit <n>", "Limit results", "50")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: { limit: string; cursor?: string }) => {
      try {
        const client = getClient();
        const results = await client.teams({
          first: parseInt(opts.limit, 10),
          after: opts.cursor,
        });
        const items = results.nodes.map((t) => ({
          id: t.id,
          name: t.name,
          key: t.key,
        }));
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
