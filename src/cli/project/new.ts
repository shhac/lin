import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveUser } from "../../lib/resolvers.ts";

export function registerNew(project: Command): void {
  project
    .command("new")
    .description("Create a new project")
    .argument("<name>", "Project name")
    .requiredOption("--team <teams>", "Team ID(s) or key(s), comma-separated")
    .option("--description <desc>", "Project description")
    .option("--lead <user>", "Project lead: name, email, or user ID")
    .option("--start-date <date>", "Start date (YYYY-MM-DD)")
    .option("--target-date <date>", "Target date (YYYY-MM-DD)")
    .option("--status <status>", "Status: backlog|planned|started|paused|completed|canceled")
    .option("--content <markdown>", "Project content body (markdown)")
    .action(
      async (
        name: string,
        opts: {
          team: string;
          description?: string;
          lead?: string;
          startDate?: string;
          targetDate?: string;
          status?: string;
          content?: string;
        },
      ) => {
        try {
          const client = getClient();

          // Validate status early
          if (opts.status) {
            const validStatuses = [
              "backlog",
              "planned",
              "started",
              "paused",
              "completed",
              "canceled",
            ];
            if (!validStatuses.includes(opts.status.toLowerCase())) {
              printError(
                `Invalid project status: "${opts.status}". Valid values: ${validStatuses.join(" | ")}`,
              );
              return;
            }
          }

          // Resolve lead name/email to user ID
          let leadId: string | undefined;
          if (opts.lead) {
            const user = await resolveUser(client, opts.lead);
            leadId = user.id;
          }

          const teamIds = opts.team.split(",").map((t) => t.trim());

          const payload = await client.createProject({
            name,
            teamIds,
            description: opts.description,
            leadId,
            startDate: opts.startDate,
            targetDate: opts.targetDate,
            state: opts.status?.toLowerCase(),
            content: opts.content,
          });
          const created = await payload.project;
          printJson({
            id: created?.id,
            slugId: created?.slugId,
            url: created?.url,
            name: created?.name,
            created: payload.success,
          });
        } catch (err) {
          printError(err instanceof Error ? err.message : "Create failed");
        }
      },
    );
}
