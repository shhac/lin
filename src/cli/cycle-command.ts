import type { Command } from "commander";
import { getClient } from "../lib/client.ts";
import { printError, printJson } from "../lib/output.ts";

export function registerCycleCommand({ program }: { program: Command }): void {
  const cycle = program.command("cycle").description("Cycle operations");

  cycle
    .command("list")
    .description("List cycles")
    .requiredOption("--team <team>", "Team ID or key")
    .option("--current", "Show only current cycle")
    .option("--next", "Show only next cycle")
    .option("--previous", "Show only previous cycle")
    .action(
      async (opts: { team: string; current?: boolean; next?: boolean; previous?: boolean }) => {
        try {
          const client = getClient();
          const team = await client.team(opts.team);
          const cycles = await team.cycles();
          printJson(
            cycles.nodes.map((c) => ({
              id: c.id,
              number: c.number,
              name: c.name,
              startsAt: c.startsAt,
              endsAt: c.endsAt,
            })),
          );
        } catch (err) {
          printError(err instanceof Error ? err.message : "List failed");
        }
      },
    );

  cycle
    .command("get")
    .description("Cycle details with issues")
    .argument("<id>", "Cycle ID")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const c = await client.cycle(id);
        const issues = await c.issues();
        printJson({
          id: c.id,
          number: c.number,
          name: c.name,
          startsAt: c.startsAt,
          endsAt: c.endsAt,
          issues: issues.nodes.map((i) => ({
            id: i.id,
            identifier: i.identifier,
            title: i.title,
            priority: i.priority,
            priorityLabel: i.priorityLabel,
          })),
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });
}
