import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printPaginated, resolvePageSize } from "../../lib/output.ts";

export function registerHistory(issue: Command): void {
  const historyCmd = issue
    .command("history")
    .description("List activity history for an issue")
    .argument("<issue-id>", "Issue ID or key")
    .option("--limit <n>", "Limit results")
    .option("--cursor <token>", "Pagination cursor for next page");
  historyCmd.action(async (issueId: string) => {
    try {
      const client = getClient();
      const opts = historyCmd.opts<{ limit?: string; cursor?: string }>();
      const i = await client.issue(issueId);
      const history = await i.history({
        first: resolvePageSize(opts),
        after: opts.cursor,
      });
      const mapped = await Promise.all(
        history.nodes.map(async (h) => {
          const [
            actor,
            fromState,
            toState,
            fromAssignee,
            toAssignee,
            fromProject,
            toProject,
            addedLabels,
            removedLabels,
          ] = await Promise.all([
            h.actor,
            h.fromState,
            h.toState,
            h.fromAssignee,
            h.toAssignee,
            h.fromProject,
            h.toProject,
            h.addedLabels ?? Promise.resolve([]),
            h.removedLabels ?? Promise.resolve([]),
          ]);
          return {
            id: h.id,
            actor: actor ? { id: actor.id, name: actor.name } : null,
            fromState: fromState ? { id: fromState.id, name: fromState.name } : null,
            toState: toState ? { id: toState.id, name: toState.name } : null,
            fromAssignee: fromAssignee ? { id: fromAssignee.id, name: fromAssignee.name } : null,
            toAssignee: toAssignee ? { id: toAssignee.id, name: toAssignee.name } : null,
            fromPriority: h.fromPriority,
            toPriority: h.toPriority,
            fromEstimate: h.fromEstimate,
            toEstimate: h.toEstimate,
            fromTitle: h.fromTitle,
            toTitle: h.toTitle,
            fromDueDate: h.fromDueDate,
            toDueDate: h.toDueDate,
            fromProject: fromProject ? { id: fromProject.id, name: fromProject.name } : null,
            toProject: toProject ? { id: toProject.id, name: toProject.name } : null,
            addedLabels: addedLabels.map((l) => ({ id: l.id, name: l.name })),
            removedLabels: removedLabels.map((l) => ({ id: l.id, name: l.name })),
            updatedDescription: h.updatedDescription,
            archived: h.archived,
            trashed: h.trashed,
            autoArchived: h.autoArchived,
            autoClosed: h.autoClosed,
            createdAt: h.createdAt,
          };
        }),
      );
      printPaginated(mapped, history.pageInfo);
    } catch (err) {
      printError(err instanceof Error ? err.message : "List history failed");
    }
  });
}
