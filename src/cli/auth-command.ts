import type { Command } from "commander";
import { LinearClient } from "@linear/sdk";
import { getApiKey, storeApiKey } from "../lib/config.ts";
import { printError, printJson } from "../lib/output.ts";

export function registerAuthCommand({ program }: { program: Command }): void {
  const auth = program.command("auth").description("Authentication management");

  auth
    .command("login")
    .description("Store Linear API key")
    .argument("<api-key>", "Linear personal API key")
    .action(async (apiKey: string) => {
      try {
        const client = new LinearClient({ apiKey });
        const viewer = await client.viewer;
        storeApiKey(apiKey);
        printJson({
          ok: true,
          user: { id: viewer.id, name: viewer.name, email: viewer.email },
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Invalid API key");
      }
    });

  auth
    .command("status")
    .description("Show current auth state and workspace info")
    .action(async () => {
      const apiKey = getApiKey();
      if (!apiKey) {
        printJson({ authenticated: false });
        return;
      }

      try {
        const client = new LinearClient({ apiKey });
        const viewer = await client.viewer;
        const org = await viewer.organization;
        printJson({
          authenticated: true,
          source: process.env.LINEAR_API_KEY ? "environment" : "config",
          user: { id: viewer.id, name: viewer.name, email: viewer.email },
          organization: { id: org.id, name: org.name, urlKey: org.urlKey },
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Auth check failed");
      }
    });
}
