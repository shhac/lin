import type { Command } from "commander";
import type { LinearDocument } from "@linear/sdk";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { mapUserSummary } from "./map-user-summary.ts";

export function registerSearch(user: Command): void {
  user
    .command("search")
    .description("Search users by name, email, or display name")
    .argument("<text>", "Search text")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (text: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const filter: LinearDocument.UserFilter = {
          or: [
            { name: { containsIgnoreCase: text } },
            { displayName: { containsIgnoreCase: text } },
            { email: { containsIgnoreCase: text } },
          ],
        };
        const results = await client.users({
          filter,
          first: resolvePageSize(opts),
          after: opts.cursor,
        });
        printPaginated(results.nodes.map(mapUserSummary), results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "Search failed");
      }
    });
}
