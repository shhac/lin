import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { PROJECT_STATUSES, PROJECT_STATUS_VALUES } from "../../lib/project-statuses.ts";
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
}
