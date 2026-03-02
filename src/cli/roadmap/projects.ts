import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { resolveRoadmap } from "../../lib/resolvers.ts";

export function registerProjects(roadmap: Command): void {
  roadmap
    .command("projects")
    .description("List projects linked to a roadmap")
    .argument("<id>", "Roadmap ID, slug, or name")
    .option("--limit <n>", "Limit results", "50")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (id: string, opts: { limit: string; cursor?: string }) => {
      try {
        const client = getClient();
        const r = await resolveRoadmap(client, id);
        const projects = await r.projects({
          first: resolvePageSize(opts),
          after: opts.cursor,
        });
        const items = await Promise.all(
          projects.nodes.map(async (p) => {
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
        printPaginated(items, projects.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get projects failed");
      }
    });
}
