import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { mapDocSummary } from "./map-doc-summary.ts";

export function registerSearch(document: Command): void {
  document
    .command("search")
    .description("Full-text search for documents")
    .argument("<text>", "Search text")
    .option("--include-comments", "Include comment text in search")
    .option("--include-archived", "Include archived documents")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (text: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const results = await client.searchDocuments(text, {
          first: resolvePageSize(opts),
          after: opts.cursor,
          includeComments: opts.includeComments !== undefined ? true : undefined,
          includeArchived: opts.includeArchived !== undefined ? true : undefined,
        });
        const items = await Promise.all(
          results.nodes.map((d) =>
            mapDocSummary(d as unknown as Parameters<typeof mapDocSummary>[0]),
          ),
        );
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "Search failed");
      }
    });
}
