import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { resolveTeam } from "../../lib/resolvers.ts";
import { mapCycleSummary } from "./map-cycle-summary.ts";

export function registerList(cycle: Command): void {
  cycle
    .command("list")
    .description("List cycles")
    .argument("<team>", "Team ID or key")
    .option("--current", "Show only current cycle")
    .option("--next", "Show only next cycle")
    .option("--previous", "Show only previous cycle")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(
      async (
        teamId: string,
        opts: {
          current?: boolean;
          next?: boolean;
          previous?: boolean;
          limit?: string;
          cursor?: string;
        },
      ) => {
        try {
          const client = getClient();
          const team = await resolveTeam(client, teamId);

          if (opts.current) {
            const c = await team.activeCycle;
            if (!c) {
              printJson([]);
              return;
            }
            printJson([mapCycleSummary(c)]);
            return;
          }

          const cycles = await team.cycles({
            first: resolvePageSize(opts),
            after: opts.cursor,
          });
          const now = new Date();

          if (opts.next) {
            const [next] = cycles.nodes
              .filter((c) => new Date(c.startsAt) > now)
              .sort((a, b) => new Date(a.startsAt).getTime() - new Date(b.startsAt).getTime());
            printJson(next ? [mapCycleSummary(next)] : []);
            return;
          }

          if (opts.previous) {
            const [prev] = cycles.nodes
              .filter((c) => new Date(c.endsAt) < now)
              .sort((a, b) => new Date(b.endsAt).getTime() - new Date(a.endsAt).getTime());
            printJson(prev ? [mapCycleSummary(prev)] : []);
            return;
          }

          printPaginated(cycles.nodes.map(mapCycleSummary), cycles.pageInfo);
        } catch (err) {
          printError(err instanceof Error ? err.message : "List failed");
        }
      },
    );
}
