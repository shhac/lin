import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { PROJECT_STATUSES, PROJECT_STATUS_VALUES } from "../../lib/project-statuses.ts";
import { PRIORITY_MAP, PRIORITY_VALUES } from "../../lib/priorities.ts";
import { resolveProject, resolveUser } from "../../lib/resolvers.ts";

export function registerUpdate(project: Command): void {
  const update = project.command("update").description("Update project fields");

  update
    .command("title")
    .description("Update project title")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<new-title>", "New title")
    .action(async (id: string, newTitle: string) => {
      try {
        const client = getClient();
        const resolved = await resolveProject(client, id);
        const payload = await client.updateProject(resolved.id, { name: newTitle });
        const p = await payload.project;
        printJson({ id: p?.id, name: p?.name, updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("status")
    .description("Update project status")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<new-status>", "New status")
    .action(async (id: string, newStatus: string) => {
      try {
        if (!PROJECT_STATUSES.includes(newStatus.toLowerCase())) {
          printError(
            `Invalid project status: "${newStatus}". Valid values: ${PROJECT_STATUS_VALUES}`,
          );
          return;
        }
        const client = getClient();
        const p = await resolveProject(client, id);
        const payload = await client.updateProject(p.id, { state: newStatus.toLowerCase() });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("description")
    .description("Update project description")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<description>", "New description")
    .action(async (id: string, description: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const payload = await client.updateProject(p.id, { description });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("content")
    .description("Update project content (markdown body)")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<content>", "New content (markdown)")
    .action(async (id: string, content: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const payload = await client.updateProject(p.id, { content });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("lead")
    .description("Update project lead")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<user>", "New lead: name, email, or user ID")
    .action(async (id: string, userInput: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const user = await resolveUser(client, userInput);
        const payload = await client.updateProject(p.id, { leadId: user.id });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("start-date")
    .description("Update project start date")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<date>", "Start date (YYYY-MM-DD)")
    .action(async (id: string, date: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const payload = await client.updateProject(p.id, { startDate: date });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("target-date")
    .description("Update project target date")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<date>", "Target date (YYYY-MM-DD)")
    .action(async (id: string, date: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const payload = await client.updateProject(p.id, { targetDate: date });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("priority")
    .description("Update project priority")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<priority>", "Priority: none|urgent|high|medium|low")
    .action(async (id: string, priority: string) => {
      try {
        const p = PRIORITY_MAP[priority.toLowerCase()];
        if (p === undefined) {
          printError(`Invalid priority: "${priority}". Valid values: ${PRIORITY_VALUES}`);
          return;
        }
        const client = getClient();
        const resolved = await resolveProject(client, id);
        const payload = await client.updateProject(resolved.id, { priority: p });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("icon")
    .description("Update project icon")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<icon>", "Icon (emoji)")
    .action(async (id: string, icon: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const payload = await client.updateProject(p.id, { icon });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("color")
    .description("Update project color")
    .argument("<id>", "Project ID, slug, or name")
    .argument("<color>", "Color (hex, e.g. #5e6ad2)")
    .action(async (id: string, color: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const payload = await client.updateProject(p.id, { color });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });
}
