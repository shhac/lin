import type { Command } from "commander";
import type { LinearDocument } from "@linear/sdk";
import { registerCommentSubcommands } from "./comment-command.ts";
import { getClient } from "../lib/client.ts";
import { buildIssueFilter } from "../lib/filters.ts";
import { printError, printJson, printPaginated } from "../lib/output.ts";

const PRIORITY_MAP: Record<string, number> = {
  none: 0,
  urgent: 1,
  high: 2,
  medium: 3,
  low: 4,
};

const PRIORITY_VALUES = "none | urgent | high | medium | low";

/** Map an issue node to a rich summary for list/search output */
async function mapIssueSummary(i: {
  id: string;
  identifier: string;
  title: string;
  branchName: string;
  priority: number;
  priorityLabel: string;
  state: Promise<{ name: string; type: string } | undefined>;
  assignee: Promise<{ id: string; name: string } | undefined>;
  team: Promise<{ key: string } | undefined>;
}): Promise<Record<string, unknown>> {
  const [state, assignee, team] = await Promise.all([i.state, i.assignee, i.team]);
  return {
    id: i.id,
    identifier: i.identifier,
    title: i.title,
    branchName: i.branchName,
    status: state ? state.name : null,
    statusType: state ? state.type : null,
    assignee: assignee ? assignee.name : null,
    assigneeId: assignee ? assignee.id : null,
    team: team ? team.key : null,
    priority: i.priority,
    priorityLabel: i.priorityLabel,
  };
}

export function registerIssueCommand({ program }: { program: Command }): void {
  const issue = program.command("issue").description("Issue operations");

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
        const results = await client.issueSearch({
          query: text,
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
    .option("--limit <n>", "Limit results", "50")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(async (opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const filter = buildIssueFilter(opts);
        const results = await client.issues({
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
        printError(err instanceof Error ? err.message : "List failed");
      }
    });

  const get = issue.command("get").description("Get issue details");

  get
    .command("overview")
    .description("Issue details: title, description, status, assignee, labels, relationships")
    .argument("<id>", "Issue ID or key (e.g. ENG-123)")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const i = await client.issue(id);
        const [assignee, state, labels, parent, team, project, comments, attachments] =
          await Promise.all([
            i.assignee,
            i.state,
            i.labels(),
            i.parent,
            i.team,
            i.project,
            i.comments(),
            i.attachments(),
          ]);
        printJson({
          id: i.id,
          identifier: i.identifier,
          url: i.url,
          title: i.title,
          description: i.description,
          branchName: i.branchName,
          status: state ? { id: state.id, name: state.name, type: state.type } : null,
          assignee: assignee ? { id: assignee.id, name: assignee.name } : null,
          team: team ? { id: team.id, key: team.key, name: team.name } : null,
          project: project ? { id: project.id, name: project.name } : null,
          priority: i.priority,
          priorityLabel: i.priorityLabel,
          commentCount: comments.nodes.length,
          labels: labels.nodes.map((l) => ({ id: l.id, name: l.name })),
          attachments: attachments.nodes.map((a) => ({
            title: a.title,
            url: a.url,
            sourceType: a.sourceType,
          })),
          parent: parent ? { id: parent.id, identifier: parent.identifier } : null,
          estimate: i.estimate,
          dueDate: i.dueDate,
          createdAt: i.createdAt,
          updatedAt: i.updatedAt,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });

  get
    .command("comments")
    .description("List comments on an issue")
    .argument("<id>", "Issue ID or key")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const i = await client.issue(id);
        const comments = await i.comments();
        const mapped = await Promise.all(
          comments.nodes.map(async (c) => {
            const user = await c.user;
            return {
              id: c.id,
              body: c.body,
              user: user ? { id: user.id, name: user.name } : null,
              createdAt: c.createdAt,
              updatedAt: c.updatedAt,
            };
          }),
        );
        printJson(mapped);
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get comments failed");
      }
    });

  issue
    .command("new")
    .description("Create a new issue")
    .argument("<title>", "Issue title")
    .requiredOption("--team <team>", "Team ID or key")
    .option("--project <project>", "Project ID, slug, or name")
    .option("--assignee <user>", "Assignee user ID")
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

        // Resolve status name to state ID
        let stateId: string | undefined;
        if (opts.status) {
          const states = await client.workflowStates();
          const state = states.nodes.find(
            (s) => s.name.toLowerCase() === opts.status!.toLowerCase(),
          );
          if (!state) {
            const validNames = [...new Set(states.nodes.map((s) => s.name))];
            printError(`Unknown status: "${opts.status}". Valid values: ${validNames.join(" | ")}`);
            return;
          }
          stateId = state.id;
        }

        const payload = await client.createIssue({
          title,
          teamId: opts.team!,
          projectId: opts.project,
          assigneeId: opts.assignee,
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
        // Look up the status ID by name
        const states = await client.workflowStates();
        const state = states.nodes.find((s) => s.name.toLowerCase() === newStatus.toLowerCase());
        if (!state) {
          const validNames = [...new Set(states.nodes.map((s) => s.name))];
          printError(`Unknown status: "${newStatus}". Valid values: ${validNames.join(" | ")}`);
          return;
        }
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
    .argument("<user-id>", "Assignee user ID")
    .action(async (id: string, userId: string) => {
      try {
        const client = getClient();
        const payload = await client.updateIssue(id, { assigneeId: userId });
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
    .argument("<project-id>", "Target project ID")
    .action(async (id: string, projectId: string) => {
      try {
        const client = getClient();
        const payload = await client.updateIssue(id, { projectId });
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

  // Comment subcommands (extracted to comment-command.ts)
  registerCommentSubcommands(issue);
}
