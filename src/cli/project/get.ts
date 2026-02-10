import type { Command } from "commander";
import type { LinearDocument } from "@linear/sdk";
import { getClient } from "../../lib/client.ts";
import { buildIssueFilter } from "../../lib/filters.ts";
import { printError, printJson, printPaginated } from "../../lib/output.ts";
import { resolveProject } from "../../lib/resolvers.ts";

export function registerGet(project: Command): void {
  const get = project.command("get").description("Get project details");

  get
    .command("overview")
    .description("Project summary: status, progress, lead, dates, milestones")
    .argument("<id>", "Project ID, slug, or name")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const lead = await p.lead;
        const milestones = await p.projectMilestones();
        printJson({
          id: p.id,
          slugId: p.slugId,
          url: p.url,
          name: p.name,
          description: p.description,
          content: p.content,
          status: p.state,
          progress: p.progress,
          lead: lead ? { id: lead.id, name: lead.name } : null,
          startDate: p.startDate,
          targetDate: p.targetDate,
          milestones: milestones.nodes.map((m) => ({
            id: m.id,
            name: m.name,
            targetDate: m.targetDate,
          })),
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });

  get
    .command("issues")
    .description("List issues within a project")
    .argument("<id>", "Project ID, slug, or name")
    .option("--status <status>", "Filter by status")
    .option("--assignee <user>", "Filter by assignee")
    .option("--priority <priority>", "Filter by priority")
    .option("--limit <n>", "Limit results", "50")
    .option("--cursor <token>", "Pagination cursor for next page")
    .action(
      async (
        id: string,
        opts: {
          status?: string;
          assignee?: string;
          priority?: string;
          limit: string;
          cursor?: string;
        },
      ) => {
        try {
          const client = getClient();
          const p = await resolveProject(client, id);
          const filter = buildIssueFilter(opts);
          const issues = await p.issues({
            first: parseInt(opts.limit, 10),
            after: opts.cursor,
            filter:
              Object.keys(filter).length > 0 ? (filter as LinearDocument.IssueFilter) : undefined,
          });
          const items = await Promise.all(
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
          printPaginated(items, issues.pageInfo);
        } catch (err) {
          printError(err instanceof Error ? err.message : "Get issues failed");
        }
      },
    );
}
