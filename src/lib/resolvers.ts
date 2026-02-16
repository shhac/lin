import type { LinearClient, Project, Roadmap, Team, User, WorkflowState } from "@linear/sdk";

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
  name: string,
  teamId: string,
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

/** Build a team filter that matches by key (e.g. "ENG") or name. */
export function buildTeamFilter(input: string): Record<string, unknown> {
  return {
    or: [{ key: { eqIgnoreCase: input } }, { name: { eqIgnoreCase: input } }],
  };
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
