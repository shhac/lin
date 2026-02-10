import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerArchive(issue: Command): void {
  issue
    .command("archive")
    .description("Archive an issue")
    .argument("<id>", "Issue ID or key")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const payload = await client.archiveIssue(id);
        printJson({ archived: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Archive failed");
      }
    });

  issue
    .command("unarchive")
    .description("Unarchive an issue")
    .argument("<id>", "Issue ID or key")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const payload = await client.unarchiveIssue(id);
        printJson({ unarchived: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Unarchive failed");
      }
    });

  issue
    .command("delete")
    .description("Delete an issue (move to trash)")
    .argument("<id>", "Issue ID or key")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const payload = await client.deleteIssue(id);
        printJson({ deleted: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Delete failed");
      }
    });
}
