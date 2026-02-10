import type { Command } from "commander";
import { getClient } from "../lib/client.ts";
import { printError, printJson } from "../lib/output.ts";

const PRIORITY_MAP: Record<string, number> = {
  none: 0,
  urgent: 1,
  high: 2,
  medium: 3,
  low: 4,
};

export function registerIssueCommand({ program }: { program: Command }): void {
  const issue = program.command("issue").description("Issue operations");

  issue
    .command("search")
    .description("Full-text search for issues")
    .argument("<text>", "Search text")
    .option("--project <project>", "Filter by project")
    .option("--team <team>", "Filter by team")
    .option("--assignee <user>", "Filter by assignee")
    .option("--status <status>", "Filter by status")
    .option("--priority <priority>", "Filter by priority")
    .option("--limit <n>", "Limit results", "50")
    .action(async (text: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const results = await client.issueSearch({
          query: text,
          first: parseInt(opts.limit ?? "50", 10),
        });
        printJson(
          results.nodes.map((i) => ({
            id: i.id,
            identifier: i.identifier,
            title: i.title,
            priority: i.priority,
            priorityLabel: i.priorityLabel,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "Search failed");
      }
    });

  issue
    .command("list")
    .description("List issues")
    .option("--project <project>", "Filter by project")
    .option("--team <team>", "Filter by team")
    .option("--assignee <user>", "Filter by assignee")
    .option("--status <status>", "Filter by status")
    .option("--priority <priority>", "Filter by priority")
    .option("--label <label>", "Filter by label")
    .option("--cycle <cycle>", "Filter by cycle")
    .option("--limit <n>", "Limit results", "50")
    .action(async (opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();
        const results = await client.issues({
          first: parseInt(opts.limit ?? "50", 10),
        });
        printJson(
          results.nodes.map((i) => ({
            id: i.id,
            identifier: i.identifier,
            title: i.title,
            priority: i.priority,
            priorityLabel: i.priorityLabel,
          })),
        );
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
        const assignee = await i.assignee;
        const state = await i.state;
        const labels = await i.labels();
        const parent = await i.parent;
        printJson({
          id: i.id,
          identifier: i.identifier,
          title: i.title,
          description: i.description,
          status: state ? { id: state.id, name: state.name, type: state.type } : null,
          assignee: assignee ? { id: assignee.id, name: assignee.name } : null,
          priority: i.priority,
          priorityLabel: i.priorityLabel,
          labels: labels.nodes.map((l) => ({ id: l.id, name: l.name })),
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
        printJson(
          comments.nodes.map((c) => ({
            id: c.id,
            body: c.body,
            createdAt: c.createdAt,
            updatedAt: c.updatedAt,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get comments failed");
      }
    });

  issue
    .command("new")
    .description("Create a new issue")
    .argument("<title>", "Issue title")
    .requiredOption("--team <team>", "Team ID or key")
    .option("--project <project>", "Project ID")
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
        const payload = await client.createIssue({
          title,
          teamId: opts.team!,
          projectId: opts.project,
          assigneeId: opts.assignee,
          priority:
            opts.priority && PRIORITY_MAP[opts.priority] !== undefined
              ? PRIORITY_MAP[opts.priority]
              : undefined,
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
          printError(`Unknown status: ${newStatus}`);
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
          printError(`Invalid priority: ${priority}. Use: none|urgent|high|medium|low`);
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

  // Comment subcommands
  const comment = issue.command("comment").description("Comment operations");

  comment
    .command("new")
    .description("Add comment to issue")
    .argument("<issue-id>", "Issue ID or key")
    .argument("<body>", "Comment body (markdown)")
    .action(async (issueId: string, body: string) => {
      try {
        const client = getClient();
        const payload = await client.createComment({ issueId, body });
        const c = await payload.comment;
        printJson({ id: c?.id, body: c?.body, created: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Comment failed");
      }
    });
}
