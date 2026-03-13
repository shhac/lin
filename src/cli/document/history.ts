import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerHistory(document: Command): void {
  document
    .command("history")
    .description("List content edit history for a document")
    .argument("<id>", "Document ID or slug ID")
    .action(async (id: string) => {
      try {
        const client = getClient();

        let docId = id;
        try {
          const d = await client.document(id);
          docId = d.id;
        } catch {
          const results = await client.documents({
            filter: { slugId: { eq: id } },
          });
          const [found] = results.nodes;
          if (!found) {
            printError(`Document not found: "${id}". Provide a UUID or slug ID.`);
            return;
          }
          docId = found.id;
        }

        const payload = await client.documentContentHistory(docId);
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
