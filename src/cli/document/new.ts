import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { resolveProject } from "../../lib/resolvers.ts";

export function registerNew(document: Command): void {
  document
    .command("new")
    .description("Create a new document")
    .argument("<title>", "Document title")
    .option("--project <project>", "Project ID, slug, or name")
    .option("--content <content>", "Document content (markdown)")
    .option("--icon <icon>", "Document icon (emoji)")
    .option("--color <color>", "Icon color (hex)")
    .action(async (title: string, opts: Record<string, string | undefined>) => {
      try {
        const client = getClient();

        let projectId: string | undefined;
        if (opts.project) {
          const resolved = await resolveProject(client, opts.project);
          projectId = resolved.id;
        }

        const payload = await client.createDocument({
          title,
          projectId,
          content: opts.content,
          icon: opts.icon,
          color: opts.color,
        });
        const created = await payload.document;
        printJson({
          id: created?.id,
          slugId: created?.slugId,
          title: created?.title,
          url: created?.url,
          created: payload.success,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Create failed");
      }
    });
}
