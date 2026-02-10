import type { Command } from "commander";
import type { LinearDocument } from "@linear/sdk";
import { getClient } from "../../lib/client.ts";
import { buildIssueFilter } from "../../lib/filters.ts";
import { printError, printPaginated } from "../../lib/output.ts";
import { mapIssueSummary } from "./map-issue-summary.ts";

export function registerSearch(issue: Command): void {
  issue
    .command("search")
    .description("Full-text search for issues")
    .argument("<text>", "Search text")
    .option("--project <project>", "Filter by project ID, slug, or name")
    .option("--team <team>", "Filter by team")
    .option("--assignee <user>", "Filter by assignee")
    .option("--status <status>", "Filter by status")
    .option("--priority <priority>", "Filter by priority")
    .option("--limit <n>", "Limit results", "50")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (text: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const filter = buildIssueFilter(opts);
        const results = await client.searchIssues(text, {
          first: parseInt(opts.limit ?? "50", 10),
          after: opts.cursor,
          filter:
            Object.keys(filter).length > 0 ? (filter as LinearDocument.IssueFilter) : undefined,
        });
        const items = await Promise.all(
          results.nodes.map((i) =>
            mapIssueSummary(i as unknown as Parameters<typeof mapIssueSummary>[0]),
          ),
        );
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "Search failed");
      }
    });
}
