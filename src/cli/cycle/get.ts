import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

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
          issues.nodes.map(async (i) => {
            const [state, assignee] = await Promise.all([i.state, i.assignee]);
            return {
              id: i.id,
              identifier: i.identifier,
              title: i.title,
              status: state ? state.name : null,
              assignee: assignee ? assignee.name : null,
              priority: i.priority,
              priorityLabel: i.priorityLabel,
            };
          }),
        );
        printJson({
          id: c.id,
          number: c.number,
          name: c.name,
          startsAt: c.startsAt,
          endsAt: c.endsAt,
          issues: mappedIssues,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });
}
