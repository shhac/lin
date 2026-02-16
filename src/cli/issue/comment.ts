import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson, printPaginated, resolvePageSize } from "../../lib/output.ts";
import { formatFileMarkdown, uploadFiles } from "../../lib/upload.ts";

function collect(val: string, prev: string[]): string[] {
  prev.push(val);
  return prev;
}

export function registerComment(issue: Command): void {
  const comment = issue.command("comment").description("Comment operations");

  const newCmd = comment
    .command("new")
    .description("Add comment to issue")
    .argument("<issue-id>", "Issue ID or key")
    .argument("<body>", "Comment body (markdown)")
    .option("--parent <comment-id>", "Parent comment ID (threaded reply)")
    .option("--file <path>", "Attach file (repeatable)", collect, []);
  newCmd.action(async (issueId: string, body: string) => {
    try {
      const client = getClient();
      const opts = newCmd.opts<{ parent?: string; file: string[] }>();
      let finalBody = body;
      if (opts.file.length > 0) {
        const uploaded = await uploadFiles(client, opts.file);
        const markdown = formatFileMarkdown(uploaded);
        finalBody = `${body}\n\n${markdown}`;
      }
      const payload = await client.createComment({
        issueId,
        body: finalBody,
        parentId: opts.parent,
      });
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
        const [user, issue, parent, children] = await Promise.all([
          c.user,
          c.issue,
          c.parent,
          c.children(),
        ]);
        printJson({
          id: c.id,
          body: c.body,
          user: user ? { id: user.id, name: user.name } : null,
          issue: issue ? { id: issue.id, identifier: issue.identifier } : null,
          parent: parent ? { id: parent.id } : null,
          childCount: children.nodes.length,
          createdAt: c.createdAt,
          updatedAt: c.updatedAt,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get comment failed");
      }
    });

  const editCmd = comment
    .command("edit")
    .description("Edit a comment")
    .argument("<comment-id>", "Comment ID")
    .argument("<body>", "New comment body (markdown)")
    .option("--file <path>", "Attach file (repeatable)", collect, []);
  editCmd.action(async (commentId: string, body: string) => {
    try {
      const client = getClient();
      const opts = editCmd.opts<{ file: string[] }>();
      let finalBody = body;
      if (opts.file.length > 0) {
        const uploaded = await uploadFiles(client, opts.file);
        finalBody = `${body}\n\n${formatFileMarkdown(uploaded)}`;
      }
      const payload = await client.updateComment(commentId, { body: finalBody });
      printJson({ updated: payload.success });
    } catch (err) {
      printError(err instanceof Error ? err.message : "Edit comment failed");
    }
  });

  const repliesCmd = comment
    .command("replies")
    .description("List replies to a comment")
    .argument("<comment-id>", "Parent comment ID")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page");
  repliesCmd.action(async (commentId: string) => {
    try {
      const client = getClient();
      const opts = repliesCmd.opts<{ limit?: string; cursor?: string }>();
      const c = await client.comment({ id: commentId });
      const children = await c.children({
        first: resolvePageSize(opts),
        after: opts.cursor,
      });
      const mapped = await Promise.all(
        children.nodes.map(async (child) => {
          const user = await child.user;
          return {
            id: child.id,
            body: child.body,
            user: user ? { id: user.id, name: user.name } : null,
            createdAt: child.createdAt,
            updatedAt: child.updatedAt,
          };
        }),
      );
      printPaginated(mapped, children.pageInfo);
    } catch (err) {
      printError(err instanceof Error ? err.message : "Get replies failed");
    }
  });
}
