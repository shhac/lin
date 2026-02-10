import type { Command } from "commander";
import { getClient } from "../lib/client.ts";
import { formatEstimateScale, getValidEstimates } from "../lib/estimates.ts";
import { printError, printJson } from "../lib/output.ts";

export function registerTeamCommand({ program }: { program: Command }): void {
  const team = program.command("team").description("Team operations");

  team
    .command("list")
    .description("List all teams")
    .action(async () => {
      try {
        const client = getClient();
        const results = await client.teams();
        printJson(
          results.nodes.map((t) => ({
            id: t.id,
            name: t.name,
            key: t.key,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });

  team
    .command("get")
    .description("Team details and members")
    .argument("<id>", "Team ID or key")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const t = await client.team(id);
        const members = await t.members();
        const estimateConfig = {
          type: t.issueEstimationType,
          allowZero: t.issueEstimationAllowZero,
          extended: t.issueEstimationExtended,
        };
        const validEstimates =
          t.issueEstimationType !== "notUsed" ? getValidEstimates(estimateConfig) : null;
        printJson({
          id: t.id,
          name: t.name,
          key: t.key,
          description: t.description,
          estimates: {
            type: t.issueEstimationType,
            allowZero: t.issueEstimationAllowZero,
            extended: t.issueEstimationExtended,
            default: t.defaultIssueEstimate,
            validValues: validEstimates,
            display: validEstimates
              ? formatEstimateScale(t.issueEstimationType, validEstimates)
              : null,
          },
          members: members.nodes.map((m) => ({
            id: m.id,
            name: m.name,
            email: m.email,
          })),
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });
}
