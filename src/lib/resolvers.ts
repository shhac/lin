import type {
  Document,
  IssueLabel,
  LinearClient,
  Project,
  Roadmap,
  Team,
  User,
  WorkflowState,
} from "@linear/sdk";
import { buildTeamFilter } from "./filters.ts";

export async function resolveUser(client: LinearClient, input: string): Promise<User> {
  const results = await client.users();
  const lower = input.toLowerCase();
  const matches = results.nodes.filter(
    (u) =>
      u.id === input ||
      u.name.toLowerCase() === lower ||
      u.email.toLowerCase() === lower ||
      u.displayName.toLowerCase() === lower,
  );
  if (matches.length === 1) {
    return matches[0]!;
  }
  if (matches.length === 0) {
    const names = results.nodes.map((u) => `${u.name} <${u.email}>`).join(", ");
    throw new Error(`User not found: "${input}". Known users: ${names}`);
  }
  const ambiguous = matches.map((u) => `${u.name} <${u.email}> (${u.id})`).join(", ");
  throw new Error(
    `Ambiguous user: "${input}" matches ${matches.length} users: ${ambiguous}. Use a unique name, email, or ID.`,
  );
}

export async function resolveDocument(client: LinearClient, input: string): Promise<Document> {
  // Try direct lookup first (works for UUIDs)
  try {
    return await client.document(input);
  } catch {
    // Fall back to search by slug ID
    const results = await client.documents({
      filter: { slugId: { eq: input } },
    });
    const [doc] = results.nodes;
    if (!doc) {
      throw new Error(`Document not found: "${input}". Provide a UUID or slug ID.`);
    }
    return doc;
  }
}

export async function resolveProject(client: LinearClient, input: string): Promise<Project> {
  // Try direct lookup first (works for UUIDs)
  try {
    return await client.project(input);
  } catch {
    // Fall back to search by slug or name
    const results = await client.projects({
      filter: {
        or: [{ slugId: { eq: input } }, { name: { eqIgnoreCase: input } }],
      },
    });
    const [project] = results.nodes;
    if (!project) {
      throw new Error(`Project not found: "${input}". Provide a UUID, slug ID, or exact name.`);
    }
    return project;
  }
}

export async function resolveWorkflowState(
  client: LinearClient,
  { name, teamId }: { name: string; teamId: string },
): Promise<WorkflowState> {
  const states = await client.workflowStates({
    filter: { team: { id: { eq: teamId } } },
  });
  const state = states.nodes.find((s) => s.name.toLowerCase() === name.toLowerCase());
  if (!state) {
    const validNames = [...new Set(states.nodes.map((s) => s.name))];
    throw new Error(`Unknown status: "${name}". Valid values: ${validNames.join(" | ")}`);
  }
  return state;
}

export async function resolveTeam(client: LinearClient, input: string): Promise<Team> {
  // Try direct lookup first (works for UUIDs)
  try {
    return await client.team(input);
  } catch {
    // Fall back to search by key or name
    const results = await client.teams({ filter: buildTeamFilter(input) });
    const [team] = results.nodes;
    if (!team) {
      const allTeams = await client.teams();
      const keys = allTeams.nodes.map((t) => `${t.key} (${t.name})`).join(", ");
      throw new Error(
        `Team not found: "${input}". Known teams: ${keys || "none"}. Provide a UUID, key, or exact name.`,
      );
    }
    return team;
  }
}

/**
 * Resolve comma-separated label names or IDs to an array of label IDs.
 * When teamId is provided, only team-scoped + workspace labels are searched,
 * reducing ambiguity from identically-named labels across teams.
 */
export async function resolveLabels(
  client: LinearClient,
  opts: { input: string; teamId?: string },
): Promise<string[]> {
  const inputs = opts.input
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
  const labels = opts.teamId
    ? await fetchTeamAndWorkspaceLabels(client, opts.teamId)
    : (await client.issueLabels()).nodes;

  const ids: string[] = [];
  for (const raw of inputs) {
    ids.push(resolveOneLabel(raw, { labels, teamScoped: !!opts.teamId }).id);
  }
  return ids;
}

async function fetchTeamAndWorkspaceLabels(
  client: LinearClient,
  teamId: string,
): Promise<IssueLabel[]> {
  const team = await client.team(teamId);
  const teamLabels = await team.labels();
  const allLabels = await client.issueLabels();
  const workspaceLabels = allLabels.nodes.filter(
    (l) => !teamLabels.nodes.some((tl) => tl.id === l.id),
  );
  return [...teamLabels.nodes, ...workspaceLabels];
}

function resolveOneLabel(
  input: string,
  ctx: { labels: IssueLabel[]; teamScoped: boolean },
): IssueLabel {
  const exact = ctx.labels.find((l) => l.id === input);
  if (exact) {
    return exact;
  }

  const lower = input.toLowerCase();
  const matches = ctx.labels.filter((l) => l.name.toLowerCase() === lower);

  if (matches.length === 1) {
    return matches[0]!;
  }

  if (matches.length === 0) {
    const names = ctx.labels.map((l) => l.name).join(", ");
    throw new Error(`Label not found: "${input}". Available labels: ${names}`);
  }

  const ambiguous = matches.map((l) => `${l.name} (${l.id})`).join(", ");
  const hint = ctx.teamScoped ? "" : " Tip: use --team to narrow scope.";
  throw new Error(
    `Ambiguous label: "${input}" matches ${matches.length} labels: ${ambiguous}. Use the label ID to disambiguate.${hint}`,
  );
}

export async function resolveRoadmap(client: LinearClient, input: string): Promise<Roadmap> {
  // Try direct lookup first (works for UUIDs)
  try {
    return await client.roadmap(input);
  } catch {
    // Fall back to search by slug or name (client.roadmaps() has no filter param)
    const results = await client.roadmaps();
    const lower = input.toLowerCase();
    const match = results.nodes.find((r) => r.slugId === input || r.name.toLowerCase() === lower);
    if (!match) {
      const names = results.nodes.map((r) => `${r.name} (${r.slugId})`).join(", ");
      throw new Error(
        `Roadmap not found: "${input}". Known roadmaps: ${names || "none"}. Provide a UUID, slug ID, or exact name.`,
      );
    }
    return match;
  }
}
