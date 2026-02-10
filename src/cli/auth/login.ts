import type { Command } from "commander";
import { LinearClient } from "@linear/sdk";
import { storeLogin, getDefaultWorkspace } from "../../lib/config.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerLogin(auth: Command): void {
  auth
    .command("login")
    .description("Store Linear API key (auto-detects workspace)")
    .argument("<api-key>", "Linear personal API key")
    .option("--alias <name>", "Custom workspace alias (default: org urlKey)")
    .action(async (apiKey: string, opts: { alias?: string }) => {
      try {
        const client = new LinearClient({ apiKey });
        const viewer = await client.viewer;
        const org = await viewer.organization;
        const alias = opts.alias ?? org.urlKey;

        storeLogin(alias, {
          api_key: apiKey,
          name: org.name,
          urlKey: org.urlKey,
        });

        const isDefault = getDefaultWorkspace() === alias;
        printJson({
          ok: true,
          user: { id: viewer.id, name: viewer.name, email: viewer.email },
          workspace: {
            alias,
            name: org.name,
            urlKey: org.urlKey,
            default: isDefault,
          },
          hint: "To add another workspace, run: lin auth login <other-api-key>",
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Invalid API key");
      }
    });
}
