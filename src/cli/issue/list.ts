import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { buildIssueFilter, nonEmptyFilter } from "../../lib/filters.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { mapIssueSummary } from "./map-issue-summary.ts";

export function registerList(issue: Command): void {
  issue
    .command("list")
    .description("List issues")
    .option("--project <project>", "Filter by project ID, slug, or name")
    .option("--team <team>", "Filter by team")
    .option("--assignee <user>", "Filter by assignee")
    .option("--status <status>", "Filter by status")
    .option("--priority <priority>", "Filter by priority")
    .option("--label <label>", "Filter by label")
    .option("--cycle <cycle>", "Filter by cycle")
    .option("--updated-after <date>", "Updated after date (YYYY-MM-DD)")
    .option("--updated-before <date>", "Updated before date (YYYY-MM-DD)")
    .option("--created-after <date>", "Created after date (YYYY-MM-DD)")
    .option("--created-before <date>", "Created before date (YYYY-MM-DD)")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const filter = buildIssueFilter(opts);
        const results = await client.issues({
          first: resolvePageSize(opts),
          after: opts.cursor,
          filter: nonEmptyFilter(filter),
        });
        const items = await Promise.all(
          results.nodes.map((i) =>
            mapIssueSummary(i as unknown as Parameters<typeof mapIssueSummary>[0]),
          ),
        );
        printPaginated(items, results.pageInfo);
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });
}
