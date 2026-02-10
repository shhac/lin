import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerList(user: Command): void {
  user
    .command("list")
    .description("List users")
    .option("--team <team>", "Filter by team")
    .action(async (opts: { team?: string }) => {
      try {
        const client = getClient();
        if (opts.team) {
          const team = await client.team(opts.team);
          const members = await team.members();
          printJson(
            members.nodes.map((u) => ({
              id: u.id,
              name: u.name,
              email: u.email,
              displayName: u.displayName,
            })),
          );
        } else {
          const results = await client.users();
          printJson(
            results.nodes.map((u) => ({
              id: u.id,
              name: u.name,
              email: u.email,
              displayName: u.displayName,
            })),
          );
        }
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
