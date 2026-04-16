import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveProject } from "../../lib/resolvers.ts";

export function registerGet(project: Command): void {
  project
    .command("get")
    .description("Project summary: status, progress, lead, dates, milestones")
    .argument("<id>", "Project ID, slug, or name")
    .action(async (id: string) => {
      try {
        const client = getClient();
        const p = await resolveProject(client, id);
        const lead = await p.lead;
        const milestones = await p.projectMilestones();
        printJson({
          id: p.id,
          slugId: p.slugId,
          url: p.url,
          name: p.name,
          description: p.description,
          content: p.content,
          status: p.state,
          progress: p.progress,
          lead: lead ? { id: lead.id, name: lead.name } : null,
          startDate: p.startDate,
          targetDate: p.targetDate,
          milestones: milestones.nodes.map((m) => ({
            id: m.id,
            name: m.name,
            targetDate: m.targetDate,
          })),
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Get failed");
      }
    });
}
