import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { resolveTeam } from "../../lib/resolvers.ts";

export function registerList(user: Command): void {
  user
    .command("list")
    .description("List users")
    .option("--team <team>", "Filter by team")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: { team?: string; limit?: string; cursor?: string }) => {
      try {
        const client = getClient();
        const first = resolvePageSize(opts);
        const after = opts.cursor;
        if (opts.team) {
          const team = await resolveTeam(client, opts.team);
          const members = await team.members({ first, after });
          const items = members.nodes.map((u) => ({
            id: u.id,
            name: u.name,
            email: u.email,
            displayName: u.displayName,
          }));
          printPaginated(items, members.pageInfo);
        } else {
          const results = await client.users({ first, after });
          const items = results.nodes.map((u) => ({
            id: u.id,
            name: u.name,
            email: u.email,
            displayName: u.displayName,
          }));
          printPaginated(items, results.pageInfo);
        }
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
