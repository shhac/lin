import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveRoadmap } from "../../lib/resolvers.ts";

export function registerGet(roadmap: Command): void {
  roadmap
    .command("get")
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
}
