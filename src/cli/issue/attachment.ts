import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerAttachment(issue: Command): void {
  const attachment = issue.command("attachment").description("Attachment operations");

  attachment
    .command("list")
    .description("List attachments on an issue")
    .argument("<issue-id>", "Issue ID or key")
    .action(async (issueId: string) => {
      try {
        const client = getClient();
        const i = await client.issue(issueId);
        const attachments = await i.attachments();
        printJson(
          attachments.nodes.map((a) => ({
            id: a.id,
            title: a.title,
            url: a.url,
            subtitle: a.subtitle,
            sourceType: a.sourceType,
          })),
        );
      } catch (err) {
        printError(err instanceof Error ? err.message : "List attachments failed");
      }
    });

  attachment
    .command("add")
    .description("Add a URL attachment to an issue")
    .argument("<issue-id>", "Issue ID or key")
    .requiredOption("--url <url>", "URL to attach")
    .requiredOption("--title <title>", "Attachment title")
    .option("--subtitle <subtitle>", "Attachment subtitle")
    .action(async (issueId: string, opts: { url: string; title: string; subtitle?: string }) => {
      try {
        const client = getClient();
        const payload = await client.createAttachment({
          issueId,
          url: opts.url,
          title: opts.title,
          subtitle: opts.subtitle,
        });
        const a = await payload.attachment;
        printJson({ id: a?.id, created: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Add attachment failed");
      }
    });

  attachment
    .command("remove")
    .description("Remove an attachment")
    .argument("<attachment-id>", "Attachment ID")
    .action(async (attachmentId: string) => {
      try {
        const client = getClient();
        const payload = await client.deleteAttachment(attachmentId);
        printJson({ deleted: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Remove attachment failed");
      }
    });
}
