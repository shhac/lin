import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { handleUnknownCommand, printError, printJson, printPaginated } from "../../lib/output.ts";
import { resolveRoadmap } from "../../lib/resolvers.ts";

export function registerGet(roadmap: Command): void {
  const get = roadmap.command("get").description("Get roadmap details");
  handleUnknownCommand(get, "Example: lin roadmap get overview <id>");

  get
    .command("overview")
    .description("Roadmap summary: name, description, owner")
    .argument("<id>", "Roadmap ID, slug, or name")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const r = await resolveRoadmap(client, id);
        const [owner, creator] = await Promise.all([r.owner, r.creator]);
        printJson({
          id: r.id,
          slugId: r.slugId,
          url: r.url,
          name: r.name,
          description: r.description,
          owner: owner ? { id: owner.id, name: owner.name } : null,
          creator: creator ? { id: creator.id, name: creator.name } : null,
          createdAt: r.createdAt,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });

  get
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
          first: parseInt(opts.limit, 10),
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
