/** Map a document node to a summary for list/search output */
export async function mapDocSummary(d: {
  id: string;
  slugId: string;
  title: string;
  url: string;
  updatedAt: Date;
  creator: Promise<{ id: string; name: string } | undefined>;
  project: Promise<{ id: string; name: string } | undefined>;
}): Promise<Record<string, unknown>> {
  const [creator, project] = await Promise.all([d.creator, d.project]);
  return {
    id: d.id,
    slugId: d.slugId,
    title: d.title,
    url: d.url,
    project: project ? { id: project.id, name: project.name } : null,
    creator: creator ? creator.name : null,
    updatedAt: d.updatedAt,
  };
}
