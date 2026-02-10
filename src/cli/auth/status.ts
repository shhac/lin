import type { Command } from "commander";
import { LinearClient } from "@linear/sdk";
import { getApiKey, getWorkspaces, getDefaultWorkspace } from "../../lib/config.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerStatus(auth: Command): void {
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
        const workspaces = getWorkspaces();
        const defaultWs = getDefaultWorkspace();
        const otherWorkspaces = Object.entries(workspaces)
          .filter(([alias]) => alias !== defaultWs)
          .map(([alias, ws]) => ({ alias, name: ws.name, urlKey: ws.urlKey }));

        printJson({
          authenticated: true,
          source: process.env.LINEAR_API_KEY ? "environment" : "config",
          user: { id: viewer.id, name: viewer.name, email: viewer.email },
          organization: { id: org.id, name: org.name, urlKey: org.urlKey },
          activeWorkspace: defaultWs,
          otherWorkspaces: otherWorkspaces.length > 0 ? otherWorkspaces : undefined,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Auth check failed");
      }
    });
}
