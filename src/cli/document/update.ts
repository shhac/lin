import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveDocument, resolveProject } from "../../lib/resolvers.ts";

export function registerUpdate(document: Command): void {
  const update = document.command("update").description("Update document fields");

  update
    .command("title")
    .description("Update document title")
    .argument("<id>", "Document ID or slug ID")
    .argument("<new-title>", "New title")
    .action(async (id: string, newTitle: string) => {
      try {
        const client = getClient();
        const doc = await resolveDocument(client, id);
        const payload = await client.updateDocument(doc.id, { title: newTitle });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("content")
    .description("Update document content")
    .argument("<id>", "Document ID or slug ID")
    .argument("<content>", "New content (markdown)")
    .action(async (id: string, content: string) => {
      try {
        const client = getClient();
        const doc = await resolveDocument(client, id);
        const payload = await client.updateDocument(doc.id, { content });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("project")
    .description("Move document to project")
    .argument("<id>", "Document ID or slug ID")
    .argument("<project>", "Project ID, slug, or name")
    .action(async (id: string, project: string) => {
      try {
        const client = getClient();
        const [doc, resolved] = await Promise.all([
          resolveDocument(client, id),
          resolveProject(client, project),
        ]);
        const payload = await client.updateDocument(doc.id, { projectId: resolved.id });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("icon")
    .description("Update document icon")
    .argument("<id>", "Document ID or slug ID")
    .argument("<icon>", "Icon (emoji)")
    .action(async (id: string, icon: string) => {
      try {
        const client = getClient();
        const doc = await resolveDocument(client, id);
        const payload = await client.updateDocument(doc.id, { icon });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });

  update
    .command("color")
    .description("Update document color")
    .argument("<id>", "Document ID or slug ID")
    .argument("<color>", "Color (hex, e.g. #5e6ad2)")
    .action(async (id: string, color: string) => {
      try {
        const client = getClient();
        const doc = await resolveDocument(client, id);
        const payload = await client.updateDocument(doc.id, { color });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Update failed");
      }
    });
}
