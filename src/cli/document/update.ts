import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveProject } from "../../lib/resolvers.ts";

export function registerUpdate(document: Command): void {
  const update = document.command("update").description("Update document fields");

  update
    .command("title")
    .description("Update document title")
    .argument("<id>", "Document ID")
    .argument("<new-title>", "New title")
    .action(async (id: string, newTitle: string) => {
      try {
        const client = getClient();
        const payload = await client.updateDocument(id, { title: newTitle });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("content")
    .description("Update document content")
    .argument("<id>", "Document ID")
    .argument("<content>", "New content (markdown)")
    .action(async (id: string, content: string) => {
      try {
        const client = getClient();
        const payload = await client.updateDocument(id, { content });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("project")
    .description("Move document to project")
    .argument("<id>", "Document ID")
    .argument("<project>", "Project ID, slug, or name")
    .action(async (id: string, project: string) => {
      try {
        const client = getClient();
        const resolved = await resolveProject(client, project);
        const payload = await client.updateDocument(id, { projectId: resolved.id });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });
}
