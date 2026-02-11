import type { Command } from "commander";
import type { LinearDocument } from "@linear/sdk";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { mapDocSummary } from "./map-doc-summary.ts";

export function registerList(document: Command): void {
  document
    .command("list")
    .description("List documents")
    .option("--project <project>", "Filter by project ID, slug, or name")
    .option("--creator <user>", "Filter by creator ID, name, or email")
    .option("--include-archived", "Include archived documents")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const filter: Record<string, unknown> = {};

        if (opts.project) {
          filter.project = {
            or: [
              { id: { eq: opts.project } },
              { slugId: { eq: opts.project } },
              { name: { eqIgnoreCase: opts.project } },
            ],
          };
        }

        if (opts.creator) {
          filter.creator = {
            or: [
              { id: { eq: opts.creator } },
              { name: { eqIgnoreCase: opts.creator } },
              { displayName: { eqIgnoreCase: opts.creator } },
              { email: { eqIgnoreCase: opts.creator } },
            ],
          };
        }

        const results = await client.documents({
          first: resolvePageSize(opts),
          after: opts.cursor,
          filter:
            Object.keys(filter).length > 0 ? (filter as LinearDocument.DocumentFilter) : undefined,
          includeArchived: opts.includeArchived !== undefined ? true : undefined,
        });
        const items = await Promise.all(
          results.nodes.map((d) =>
            mapDocSummary(d as unknown as Parameters<typeof mapDocSummary>[0]),
          ),
        );
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
