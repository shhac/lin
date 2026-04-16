import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { mapIssueSummary } from "../issue/map-issue-summary.ts";
import { mapCycleSummary } from "./map-cycle-summary.ts";

export function registerGet(cycle: Command): void {
  cycle
    .command("get")
    .description("Cycle details with issues")
    .argument("<id>", "Cycle ID")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const c = await client.cycle(id);
        const issues = await c.issues();
        const mappedIssues = await Promise.all(
          issues.nodes.map((i) =>
            mapIssueSummary(i as unknown as Parameters<typeof mapIssueSummary>[0]),
          ),
        );
        printJson({
          ...mapCycleSummary(c),
          issues: mappedIssues,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });
}
