import type { Command } from "commander";
import type { LinearDocument } from "@linear/sdk";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { buildTeamFilter } from "../../lib/resolvers.ts";

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
        const filter: LinearDocument.ProjectFilter = {};
        if (opts.team) {
          filter.accessibleTeams = { some: buildTeamFilter(opts.team) };
        }
        if (opts.status) {
          filter.state = { eqIgnoreCase: opts.status };
        }
        const results = await client.projects({
          first: resolvePageSize(opts),
          after: opts.cursor,
          filter: Object.keys(filter).length > 0 ? filter : undefined,
        });
        const items = await Promise.all(
          results.nodes.map(async (p) => {
            const lead = await p.lead;
            return {
              id: p.id,
              slugId: p.slugId,
              url: p.url,
              name: p.name,
              status: p.state,
              progress: p.progress,
              lead: lead ? lead.name : null,
              startDate: p.startDate,
              targetDate: p.targetDate,
            };
          }),
        );
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
