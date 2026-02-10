import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerStates(team: Command): void {
  team
    .command("states")
    .description("List workflow states for a team (discover valid status values)")
    .argument("<team>", "Team key or name")
    .action(async (teamInput: string) => {
      try {
        const client = getClient();
        const states = await client.workflowStates({
          filter: {
            team: {
              or: [{ key: { eqIgnoreCase: teamInput } }, { name: { eqIgnoreCase: teamInput } }],
            },
          },
        });
        if (states.nodes.length === 0) {
          printError(`No workflow states found for team "${teamInput}".`);
          return;
        }
        const sorted = [...states.nodes].sort((a, b) => a.position - b.position);
        printJson(
          sorted.map((s) => ({
            id: s.id,
            name: s.name,
            type: s.type,
            color: s.color,
            position: s.position,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "Failed to list workflow states");
      }
    });
}
