import { resolvePriority } from "./priorities.ts";

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/** Return the filter object if non-empty, otherwise undefined. */
export function nonEmptyFilter<T>(filter: Record<string, unknown>): T | undefined {
  return Object.keys(filter).length > 0 ? (filter as T) : undefined;
}

/** Build a team filter that matches by key (e.g. "ENG") or name. */
export function buildTeamFilter(input: string): Record<string, unknown> {
  return {
    or: [{ key: { eqIgnoreCase: input } }, { name: { eqIgnoreCase: input } }],
  };
}

/** Build a project filter that matches by ID, slugId, or name. */
export function buildProjectFilter(input: string): Record<string, unknown> {
  const branches: Record<string, unknown>[] = [
    { slugId: { eq: input } },
    { name: { eqIgnoreCase: input } },
  ];
  if (UUID_RE.test(input)) {
    branches.unshift({ id: { eq: input } });
  }
  return { or: branches };
}

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
    filter.project = buildProjectFilter(opts.project);
  }

  if (opts.team) {
    filter.team = buildTeamFilter(opts.team);
  }

  if (opts.assignee) {
    if (opts.assignee.toLowerCase() === "me") {
      filter.assignee = { isMe: { eq: true } };
    } else {
      const branches: Record<string, unknown>[] = [
        { name: { eqIgnoreCase: opts.assignee } },
        { displayName: { eqIgnoreCase: opts.assignee } },
        { email: { eqIgnoreCase: opts.assignee } },
      ];
      if (UUID_RE.test(opts.assignee)) {
        branches.unshift({ id: { eq: opts.assignee } });
      }
      filter.assignee = { or: branches };
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
