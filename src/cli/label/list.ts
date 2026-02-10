import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated } from "../../lib/output.ts";
import { resolveTeam } from "../../lib/resolvers.ts";

export function registerList(label: Command): void {
  label
    .command("list")
    .description("List labels")
    .option("--team <team>", "Filter by team")
    .option("--limit <n>", "Limit results", "50")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: { team?: string; limit: string; cursor?: string }) => {
      try {
        const client = getClient();
        const first = parseInt(opts.limit, 10);
        const after = opts.cursor;
        if (opts.team) {
          const team = await resolveTeam(client, opts.team);
          const labels = await team.labels({ first, after });
          const items = labels.nodes.map((l) => ({
            id: l.id,
            name: l.name,
            color: l.color,
          }));
          printPaginated(items, labels.pageInfo);
        } else {
          const results = await client.issueLabels({ first, after });
          const items = results.nodes.map((l) => ({
            id: l.id,
            name: l.name,
            color: l.color,
          }));
          printPaginated(items, results.pageInfo);
        }
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
