/** Map a project node to a summary for list/search output */
export async function mapProjectSummary(p: {
  id: string;
  slugId: string;
  url: string;
  name: string;
  state: string;
  progress: number;
  startDate?: string;
  targetDate?: string;
  lead: Promise<{ name: string } | undefined>;
}): Promise<Record<string, unknown>> {
  const lead = await p.lead;
  return {
    id: p.id,
    slugId: p.slugId,
    url: p.url,
    name: p.name,
    status: p.state,
    progress: p.progress,
    lead: lead ? lead.name : null,
    startDate: p.startDate,
    targetDate: p.targetDate,
  };
}
