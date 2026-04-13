import { resolvePriority } from "./priorities.ts";
import { buildTeamFilter } from "./resolvers.ts";

/**
 * Build an IssueFilter from CLI options.
 * Accepts the common filter flags: --project, --team, --assignee, --status, --priority, --label, --cycle.
 * Returns a partial filter object suitable for passing to the Linear SDK.
 */
export function buildIssueFilter(
  opts: Record<string, string | undefined>,
): Record<string, unknown> {
  const filter: Record<string, unknown> = {};

  if (opts.project) {
    filter.project = {
      or: [
        { id: { eq: opts.project } },
        { slugId: { eq: opts.project } },
        { name: { eqIgnoreCase: opts.project } },
      ],
    };
  }

  if (opts.team) {
    filter.team = buildTeamFilter(opts.team);
  }

  if (opts.assignee) {
    if (opts.assignee.toLowerCase() === "me") {
      filter.assignee = { isMe: { eq: true } };
    } else {
      // Accept user ID, name, display name, or email
      filter.assignee = {
        or: [
          { id: { eq: opts.assignee } },
          { name: { eqIgnoreCase: opts.assignee } },
          { displayName: { eqIgnoreCase: opts.assignee } },
          { email: { eqIgnoreCase: opts.assignee } },
        ],
      };
    }
  }

  if (opts.status) {
    filter.state = { name: { eqIgnoreCase: opts.status } };
  }

  if (opts.priority) {
    const p = resolvePriority(opts.priority);
    if (p !== undefined) {
      filter.priority = { eq: p };
    }
  }

  if (opts.label) {
    filter.labels = { name: { eqIgnoreCase: opts.label } };
  }

  if (opts.cycle) {
    filter.cycle = { id: { eq: opts.cycle } };
  }

  const updatedAt: Record<string, string> = {};
  if (opts["updated-after"]) {
    updatedAt.gte = opts["updated-after"];
  }
  if (opts["updated-before"]) {
    updatedAt.lte = opts["updated-before"];
  }
  if (Object.keys(updatedAt).length > 0) {
    filter.updatedAt = updatedAt;
  }

  const createdAt: Record<string, string> = {};
  if (opts["created-after"]) {
    createdAt.gte = opts["created-after"];
  }
  if (opts["created-before"]) {
    createdAt.lte = opts["created-before"];
  }
  if (Object.keys(createdAt).length > 0) {
    filter.createdAt = createdAt;
  }

  return filter;
}
