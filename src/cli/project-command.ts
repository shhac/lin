import type { Command } from "commander";
import { getClient } from "../lib/client.ts";
import { printError, printJson } from "../lib/output.ts";

export function registerProjectCommand({ program }: { program: Command }): void {
  const project = program.command("project").description("Project operations");

  project
    .command("search")
    .description("Search projects by name/description")
    .argument("<text>", "Search text")
    .action(async (text: string) => {
      try {
        const client = getClient();
        const results = await client.projects({
          filter: { name: { containsIgnoreCase: text } },
        });
        printJson(
          results.nodes.map((p) => ({
            id: p.id,
            name: p.name,
            status: p.state,
            progress: p.progress,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "Search failed");
      }
    });

  project
    .command("list")
    .description("List all projects")
    .option("--team <team>", "Filter by team name or ID")
    .option("--status <status>", "Filter by status")
    .option("--limit <n>", "Limit results", "50")
    .action(async (opts: { team?: string; status?: string; limit: string }) => {
      try {
        const client = getClient();
        const results = await client.projects({
          first: parseInt(opts.limit, 10),
        });
        printJson(
          results.nodes.map((p) => ({
            id: p.id,
            name: p.name,
            status: p.state,
            progress: p.progress,
            startDate: p.startDate,
            targetDate: p.targetDate,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "List failed");
      }
    });

  const get = project.command("get").description("Get project details");

  get
    .command("overview")
    .description("Project details: status, progress, lead, dates, milestones")
    .argument("<id>", "Project ID or name")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const p = await client.project(id);
        const lead = await p.lead;
        const milestones = await p.projectMilestones();
        printJson({
          id: p.id,
          name: p.name,
          description: p.description,
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
    .argument("<id>", "Project ID or name")
    .option("--status <status>", "Filter by status")
    .option("--assignee <user>", "Filter by assignee")
    .option("--priority <priority>", "Filter by priority")
    .option("--limit <n>", "Limit results", "50")
    .action(
      async (
        id: string,
        opts: { status?: string; assignee?: string; priority?: string; limit: string },
      ) => {
        try {
          const client = getClient();
          const p = await client.project(id);
          const issues = await p.issues({ first: parseInt(opts.limit, 10) });
          printJson(
            issues.nodes.map((i) => ({
              id: i.id,
              identifier: i.identifier,
              title: i.title,
              priority: i.priority,
              priorityLabel: i.priorityLabel,
            })),
          );
        } catch (err) {
          printError(err instanceof Error ? err.message : "Get issues failed");
        }
      },
    );

  const update = project.command("update").description("Update project fields");

  update
    .command("title")
    .description("Update project title")
    .argument("<id>", "Project ID")
    .argument("<new-title>", "New title")
    .action(async (id: string, newTitle: string) => {
      try {
        const client = getClient();
        const payload = await client.updateProject(id, { name: newTitle });
        const p = await payload.project;
        printJson({ id: p?.id, name: p?.name, updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("status")
    .description("Update project status")
    .argument("<id>", "Project ID")
    .argument("<new-status>", "New status")
    .action(async (id: string, newStatus: string) => {
      try {
        const client = getClient();
        const payload = await client.updateProject(id, { state: newStatus });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("description")
    .description("Update project description")
    .argument("<id>", "Project ID")
    .argument("<description>", "New description")
    .action(async (id: string, description: string) => {
      try {
        const client = getClient();
        const payload = await client.updateProject(id, { description });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("lead")
    .description("Update project lead")
    .argument("<id>", "Project ID")
    .argument("<user-id>", "New lead user ID")
    .action(async (id: string, userId: string) => {
      try {
        const client = getClient();
        const payload = await client.updateProject(id, { leadId: userId });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });
}
