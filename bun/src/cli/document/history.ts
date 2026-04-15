import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveDocument } from "../../lib/resolvers.ts";

export function registerHistory(document: Command): void {
  document
    .command("history")
    .description("List content edit history for a document")
    .argument("<id>", "Document ID or slug ID")
    .action(async (id: string) => {
      try {
        const client = getClient();

        const doc = await resolveDocument(client, id);
        const payload = await client.documentContentHistory(doc.id);
        const items = payload.history.map((h) => ({
          id: h.id,
          actorIds: h.actorIds,
          contentDataSnapshotAt: h.contentDataSnapshotAt,
          createdAt: h.createdAt,
        }));
        printJson({ items });
      } catch (err) {
        printError(err instanceof Error ? err.message : "List history failed");
      }
    });
}
