import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerList(team: Command): void {
  team
    .command("list")
    .description("List all teams")
    .action(async () => {
      try {
        const client = getClient();
        const results = await client.teams();
        printJson(
          results.nodes.map((t) => ({
            id: t.id,
            name: t.name,
            key: t.key,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
