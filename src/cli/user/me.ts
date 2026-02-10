import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerMe(user: Command): void {
  user
    .command("me")
    .description("Current authenticated user")
    .action(async () => {
      try {
        const client = getClient();
        const viewer = await client.viewer;
        const org = await viewer.organization;
        printJson({
          id: viewer.id,
          name: viewer.name,
          email: viewer.email,
          displayName: viewer.displayName,
          organization: { id: org.id, name: org.name },
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Failed to get user");
      }
    });
}
