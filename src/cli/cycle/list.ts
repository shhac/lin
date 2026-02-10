import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveTeam } from "../../lib/resolvers.ts";

export function registerList(cycle: Command): void {
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
          const team = await resolveTeam(client, opts.team);

          if (opts.current) {
            const c = await team.activeCycle;
            if (!c) {
              printJson([]);
              return;
            }
            printJson([
              {
                id: c.id,
                number: c.number,
                name: c.name,
                startsAt: c.startsAt,
                endsAt: c.endsAt,
              },
            ]);
            return;
          }

          const cycles = await team.cycles();
          const now = new Date();

          if (opts.next) {
            const [next] = cycles.nodes
              .filter((c) => new Date(c.startsAt) > now)
              .sort((a, b) => new Date(a.startsAt).getTime() - new Date(b.startsAt).getTime());
            printJson(
              next
                ? [
                    {
                      id: next.id,
                      number: next.number,
                      name: next.name,
                      startsAt: next.startsAt,
                      endsAt: next.endsAt,
                    },
                  ]
                : [],
            );
            return;
          }

          if (opts.previous) {
            const [prev] = cycles.nodes
              .filter((c) => new Date(c.endsAt) < now)
              .sort((a, b) => new Date(b.endsAt).getTime() - new Date(a.endsAt).getTime());
            printJson(
              prev
                ? [
                    {
                      id: prev.id,
                      number: prev.number,
                      name: prev.name,
                      startsAt: prev.startsAt,
                      endsAt: prev.endsAt,
                    },
                  ]
                : [],
            );
            return;
          }

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
}
