import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

const VALID_TYPES = ["blocks", "duplicate", "related"];
const VALID_TYPES_DISPLAY = VALID_TYPES.join(" | ");

export function registerRelation(issue: Command): void {
  const relation = issue.command("relation").description("Issue relation operations");

  relation
    .command("list")
    .description("List all relations on an issue (both directions)")
    .argument("<issue-id>", "Issue ID or key")
    .action(async (issueId: string) => {
      try {
        const client = getClient();
        const i = await client.issue(issueId);
        const [relations, inverseRelations] = await Promise.all([
          i.relations(),
          i.inverseRelations(),
        ]);

        const mapped = await Promise.all([
          ...relations.nodes.map(async (r) => {
            const related = await r.relatedIssue;
            return {
              id: r.id,
              type: r.type,
              relatedIssue: related?.identifier,
            };
          }),
          ...inverseRelations.nodes.map(async (r) => {
            const source = await r.issue;
            return {
              id: r.id,
              type: r.type === "blocks" ? "blocked_by" : r.type,
              relatedIssue: source?.identifier,
            };
          }),
        ]);

        printJson(mapped);
      } catch (err) {
        printError(err instanceof Error ? err.message : "List relations failed");
      }
    });

  relation
    .command("add")
    .description("Add a relation between two issues")
    .argument("<issue-id>", "Source issue ID or key")
    .requiredOption("--type <type>", `Relation type: ${VALID_TYPES_DISPLAY}`)
    .requiredOption("--related <related-id>", "Target issue ID or key")
    .action(async (issueId: string, opts: { type: string; related: string }) => {
      try {
        const normalized = opts.type.toLowerCase();
        if (!VALID_TYPES.includes(normalized)) {
          printError(
            `Invalid relation type: "${opts.type}". Valid values: ${VALID_TYPES_DISPLAY}`,
          );
          return;
        }
        const client = getClient();
        // IssueRelationType enum isn't re-exported from @linear/sdk; values match after validation
        const payload = await client.createIssueRelation({
          issueId,
          relatedIssueId: opts.related,
          // eslint-disable-next-line @typescript-eslint/no-explicit-any -- enum not exported
          type: normalized as any,
        });
        printJson({ created: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Add relation failed");
      }
    });

  relation
    .command("remove")
    .description("Remove a relation")
    .argument("<relation-id>", "Relation ID")
    .action(async (relationId: string) => {
      try {
        const client = getClient();
        const payload = await client.deleteIssueRelation(relationId);
        printJson({ deleted: payload.success });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Remove relation failed");
      }
    });
}
