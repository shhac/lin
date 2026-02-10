import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated } from "../../lib/output.ts";

export function registerSearch(project: Command): void {
  project
    .command("search")
    .description("Search projects by name/description")
    .argument("<text>", "Search text")
    .action(async (text: string) => {
      try {
        const client = getClient();
        const results = await client.projects({
          filter: { name: { containsIgnoreCase: text } },
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
