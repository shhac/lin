import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerGet(issue: Command): void {
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
}
