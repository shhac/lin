import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { PRIORITY_VALUES, resolvePriority } from "../../lib/priorities.ts";
import {
  resolveLabels,
  resolveProject,
  resolveTeam,
  resolveUser,
  resolveWorkflowState,
} from "../../lib/resolvers.ts";

export function registerNew(issue: Command): void {
  issue
    .command("new")
    .description("Create a new issue")
    .argument("<title>", "Issue title")
    .requiredOption("--team <team>", "Team ID or key")
    .option("--project <project>", "Project ID, slug, or name")
    .option("--assignee <user>", "Assignee: name, email, or user ID")
    .option("--priority <priority>", "Priority: none|urgent|high|medium|low")
    .option("--status <status>", "Status name")
    .option("--labels <labels>", "Comma-separated label names or IDs")
    .option("--description <desc>", "Issue description (markdown)")
    .option("--cycle <cycle>", "Cycle ID")
    .option("--parent <parent>", "Parent issue ID")
    .option("--estimate <points>", "Estimate points")
    .action(async (title: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();

        const priority = opts.priority ? resolvePriority(opts.priority) : undefined;
        if (opts.priority && priority === undefined) {
          printError(`Invalid priority: "${opts.priority}". Valid values: ${PRIORITY_VALUES}`);
          return;
        }

        const team = await resolveTeam(client, opts.team!);
        const stateId = opts.status
          ? (await resolveWorkflowState(client, { name: opts.status, teamId: team.id })).id
          : undefined;
        const assigneeId = opts.assignee
          ? (await resolveUser(client, opts.assignee)).id
          : undefined;
        const projectId = opts.project
          ? (await resolveProject(client, opts.project)).id
          : undefined;
        const labelIds = opts.labels
          ? await resolveLabels(client, { input: opts.labels, teamId: team.id })
          : undefined;

        const payload = await client.createIssue({
          title,
          teamId: team.id,
          projectId,
          assigneeId,
          stateId,
          priority,
          description: opts.description,
          cycleId: opts.cycle,
          parentId: opts.parent,
          estimate: opts.estimate ? parseInt(opts.estimate, 10) : undefined,
          labelIds,
        });
        const created = await payload.issue;
        printJson({
          id: created?.id,
          identifier: created?.identifier,
          url: created?.url,
          title: created?.title,
          created: payload.success,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Create failed");
      }
    });
}
