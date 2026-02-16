import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { formatEstimateScale, getValidEstimates } from "../../lib/estimates.ts";
import { printError, printJson } from "../../lib/output.ts";
import { PRIORITY_MAP, PRIORITY_VALUES } from "../../lib/priorities.ts";
import { resolveProject, resolveUser, resolveWorkflowState } from "../../lib/resolvers.ts";

export function registerUpdate(issue: Command): void {
  const update = issue.command("update").description("Update issue fields");

  update
    .command("title")
    .description("Update issue title")
    .argument("<id>", "Issue ID or key")
    .argument("<new-title>", "New title")
    .action(async (id: string, newTitle: string) => {
      try {
        const client = getClient();
        const payload = await client.updateIssue(id, { title: newTitle });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("status")
    .description("Update issue status")
    .argument("<id>", "Issue ID or key")
    .argument("<new-status>", "New status name")
    .action(async (id: string, newStatus: string) => {
      try {
        const client = getClient();
        const issue = await client.issue(id);
        const team = await issue.team;
        if (!team) {
          printError("Could not resolve team for this issue.");
          return;
        }
        const state = await resolveWorkflowState(client, newStatus, team.id);
        const payload = await client.updateIssue(id, { stateId: state.id });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("assignee")
    .description("Update issue assignee")
    .argument("<id>", "Issue ID or key")
    .argument("<user>", "Assignee: name, email, or user ID")
    .action(async (id: string, userInput: string) => {
      try {
        const client = getClient();
        const user = await resolveUser(client, userInput);
        const payload = await client.updateIssue(id, { assigneeId: user.id });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("priority")
    .description("Update issue priority")
    .argument("<id>", "Issue ID or key")
    .argument("<priority>", "Priority: none|urgent|high|medium|low")
    .action(async (id: string, priority: string) => {
      try {
        const p = PRIORITY_MAP[priority.toLowerCase()];
        if (p === undefined) {
          printError(`Invalid priority: "${priority}". Valid values: ${PRIORITY_VALUES}`);
          return;
        }
        const client = getClient();
        const payload = await client.updateIssue(id, { priority: p });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("project")
    .description("Move issue to project")
    .argument("<id>", "Issue ID or key")
    .argument("<project>", "Project ID, slug, or name")
    .action(async (id: string, project: string) => {
      try {
        const client = getClient();
        const resolved = await resolveProject(client, project);
        const payload = await client.updateIssue(id, { projectId: resolved.id });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("labels")
    .description("Set issue labels")
    .argument("<id>", "Issue ID or key")
    .argument("<labels>", "Comma-separated label IDs")
    .action(async (id: string, labels: string) => {
      try {
        const client = getClient();
        const payload = await client.updateIssue(id, { labelIds: labels.split(",") });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("description")
    .description("Update issue description")
    .argument("<id>", "Issue ID or key")
    .argument("<description>", "New description (markdown)")
    .action(async (id: string, description: string) => {
      try {
        const client = getClient();
        const payload = await client.updateIssue(id, { description });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("estimate")
    .description("Update issue estimate (validated against team scale)")
    .argument("<id>", "Issue ID or key")
    .argument("<value>", "Estimate value (number)")
    .action(async (id: string, value: string) => {
      try {
        const estimate = parseInt(value, 10);
        if (Number.isNaN(estimate)) {
          printError(`Invalid estimate: "${value}". Must be a number.`);
          return;
        }
        const client = getClient();
        const i = await client.issue(id);
        const team = await i.team;
        if (!team) {
          printError("Could not resolve team for this issue.");
          return;
        }
        if (team.issueEstimationType === "notUsed") {
          printError(`Team "${team.key}" does not use estimates.`);
          return;
        }
        const estimateConfig = {
          type: team.issueEstimationType,
          allowZero: team.issueEstimationAllowZero,
          extended: team.issueEstimationExtended,
        };
        const valid = getValidEstimates(estimateConfig);
        if (!valid.includes(estimate)) {
          const scale = formatEstimateScale(team.issueEstimationType, valid);
          printError(
            `Invalid estimate: ${estimate}. Team "${team.key}" uses ${team.issueEstimationType} scale. Valid values: ${scale}`,
          );
          return;
        }
        const payload = await client.updateIssue(id, { estimate });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });
}
