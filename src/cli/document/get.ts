import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveDocument } from "../../lib/resolvers.ts";

export function registerGet(document: Command): void {
  document
    .command("get")
    .description("Get document details (includes full content)")
    .argument("<id>", "Document ID or slug ID")
    .action(async (id: string) => {
      try {
        const client = getClient();

        const d = await resolveDocument(client, id);

        const [creator, project, updatedBy] = await Promise.all([
          d.creator,
          d.project,
          d.updatedBy,
        ]);
        printJson({
          id: d.id,
          slugId: d.slugId,
          title: d.title,
          content: d.content,
          url: d.url,
          icon: d.icon,
          color: d.color,
          project: project ? { id: project.id, name: project.name, slugId: project.slugId } : null,
          creator: creator ? { id: creator.id, name: creator.name } : null,
          updatedBy: updatedBy ? { id: updatedBy.id, name: updatedBy.name } : null,
          createdAt: d.createdAt,
          updatedAt: d.updatedAt,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });
}
