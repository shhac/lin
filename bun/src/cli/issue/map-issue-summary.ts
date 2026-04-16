/** Map an issue node to a rich summary for list/search output */
export async function mapIssueSummary(i: {
  id: string;
  identifier: string;
  title: string;
  branchName: string;
  priority: number;
  priorityLabel: string;
  state: Promise<{ name: string; type: string } | undefined>;
  assignee: Promise<{ id: string; name: string } | undefined>;
  team: Promise<{ key: string } | undefined>;
}): Promise<Record<string, unknown>> {
  const [state, assignee, team] = await Promise.all([i.state, i.assignee, i.team]);
  return {
    id: i.id,
    identifier: i.identifier,
    title: i.title,
    branchName: i.branchName,
    status: state ? state.name : null,
    statusType: state ? state.type : null,
    assignee: assignee ? assignee.name : null,
    assigneeId: assignee ? assignee.id : null,
    team: team ? team.key : null,
    priority: i.priority,
    priorityLabel: i.priorityLabel,
  };
}
