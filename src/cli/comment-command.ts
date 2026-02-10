import type { Command } from "commander";
import { getClient } from "../lib/client.ts";
import { printError, printJson } from "../lib/output.ts";

export function registerCommentSubcommands(issue: Command): void {
  const comment = issue.command("comment").description("Comment operations");

  comment
    .command("new")
    .description("Add comment to issue")
    .argument("<issue-id>", "Issue ID or key")
    .argument("<body>", "Comment body (markdown)")
    .action(async (issueId: string, body: string) => {
      try {
        const client = getClient();
        const payload = await client.createComment({ issueId, body });
        const c = await payload.comment;
        printJson({ id: c?.id, body: c?.body, created: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Comment failed");
      }
    });

  comment
    .command("get")
    .description("Get a specific comment")
    .argument("<comment-id>", "Comment ID")
    .action(async (commentId: string) => {
      try {
        const client = getClient();
        const c = await client.comment({ id: commentId });
        const user = await c.user;
        const issue = await c.issue;
        printJson({
          id: c.id,
          body: c.body,
          user: user ? { id: user.id, name: user.name } : null,
          issue: issue ? { id: issue.id, identifier: issue.identifier } : null,
          createdAt: c.createdAt,
          updatedAt: c.updatedAt,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get comment failed");
      }
    });

  comment
    .command("edit")
    .description("Edit a comment")
    .argument("<comment-id>", "Comment ID")
    .argument("<body>", "New comment body (markdown)")
    .action(async (commentId: string, body: string) => {
      try {
        const client = getClient();
        const payload = await client.updateComment(commentId, { body });
        printJson({ updated: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Edit comment failed");
      }
    });
}
