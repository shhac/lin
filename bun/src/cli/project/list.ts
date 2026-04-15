import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { buildTeamFilter, nonEmptyFilter } from "../../lib/filters.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { mapProjectSummary } from "./map-project-summary.ts";

export function registerList(project: Command): void {
  project
    .command("list")
    .description("List all projects")
    .option("--team <team>", "Filter by team name or ID")
    .option("--status <status>", "Filter by status")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: { team?: string; status?: string; limit?: string; cursor?: string }) => {
      try {
        const client = getClient();
        const filter: Record<string, unknown> = {};
        if (opts.team) {
          filter.accessibleTeams = { some: buildTeamFilter(opts.team) };
        }
        if (opts.status) {
          filter.state = { eqIgnoreCase: opts.status };
        }
        const results = await client.projects({
          first: resolvePageSize(opts),
          after: opts.cursor,
          filter: nonEmptyFilter(filter),
        });
        const items = await Promise.all(
          results.nodes.map((p) =>
            mapProjectSummary(p as unknown as Parameters<typeof mapProjectSummary>[0]),
          ),
        );
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
