const PRIORITY_MAP: Record<string, number> = {
  none: 0,
  urgent: 1,
  high: 2,
  medium: 3,
  low: 4,
};

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
    filter.project = { id: { eq: opts.project } };
  }

  if (opts.team) {
    // Accept team key (e.g. "ENG") or team name; try key first via OR filter
    filter.team = {
      or: [{ key: { eqIgnoreCase: opts.team } }, { name: { eqIgnoreCase: opts.team } }],
    };
  }

  if (opts.assignee) {
    // Accept user ID, name, or display name
    filter.assignee = {
      or: [
        { id: { eq: opts.assignee } },
        { name: { eqIgnoreCase: opts.assignee } },
        { displayName: { eqIgnoreCase: opts.assignee } },
      ],
    };
  }

  if (opts.status) {
    filter.state = { name: { eqIgnoreCase: opts.status } };
  }

  if (opts.priority) {
    const p = PRIORITY_MAP[opts.priority.toLowerCase()];
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

  return filter;
}
