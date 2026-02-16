import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { PRIORITY_MAP, PRIORITY_VALUES } from "../../lib/priorities.ts";
import {
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
    .option("--labels <labels>", "Comma-separated label IDs")
    .option("--description <desc>", "Issue description (markdown)")
    .option("--cycle <cycle>", "Cycle ID")
    .option("--parent <parent>", "Parent issue ID")
    .option("--estimate <points>", "Estimate points")
    .action(async (title: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();

        // Validate priority early
        let priority: number | undefined;
        if (opts.priority) {
          priority = PRIORITY_MAP[opts.priority.toLowerCase()];
          if (priority === undefined) {
            printError(`Invalid priority: "${opts.priority}". Valid values: ${PRIORITY_VALUES}`);
            return;
          }
        }

        // Resolve team key/name/UUID to team object
        const team = await resolveTeam(client, opts.team!);

        // Resolve status name to state ID (scoped to team)
        let stateId: string | undefined;
        if (opts.status) {
          const state = await resolveWorkflowState(client, opts.status, team.id);
          stateId = state.id;
        }

        // Resolve assignee name/email to user ID
        let assigneeId: string | undefined;
        if (opts.assignee) {
          const user = await resolveUser(client, opts.assignee);
          assigneeId = user.id;
        }

        // Resolve project name/slug/UUID to project ID
        let projectId: string | undefined;
        if (opts.project) {
          const project = await resolveProject(client, opts.project);
          projectId = project.id;
        }

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
          labelIds: opts.labels ? opts.labels.split(",") : undefined,
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
