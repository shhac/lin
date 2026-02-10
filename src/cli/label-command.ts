import type { Command } from "commander";
import { getClient } from "../lib/client.ts";
import { printError, printJson } from "../lib/output.ts";

export function registerLabelCommand({ program }: { program: Command }): void {
  const label = program.command("label").description("Label operations");

  label
    .command("list")
    .description("List labels")
    .option("--team <team>", "Filter by team")
    .action(async (opts: { team?: string }) => {
      try {
        const client = getClient();
        if (opts.team) {
          const team = await client.team(opts.team);
          const labels = await team.labels();
          printJson(
            labels.nodes.map((l) => ({
              id: l.id,
              name: l.name,
              color: l.color,
            })),
          );
        } else {
          const results = await client.issueLabels();
          printJson(
            results.nodes.map((l) => ({
              id: l.id,
              name: l.name,
              color: l.color,
            })),
          );
        }
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
