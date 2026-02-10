import type { LinearClient, Project } from "@linear/sdk";

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
