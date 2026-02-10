import type { Command } from "commander";
import { LinearClient } from "@linear/sdk";
import {
  getApiKey,
  clearApiKey,
  clearAll,
  storeLogin,
  getWorkspaces,
  getDefaultWorkspace,
  setDefaultWorkspace,
  removeWorkspace,
} from "../lib/config.ts";
import { printError, printJson } from "../lib/output.ts";

export function registerAuthCommand({ program }: { program: Command }): void {
  const auth = program.command("auth").description("Authentication management");

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

  auth
    .command("logout")
    .description("Clear stored credentials")
    .option("--all", "Remove all workspaces (default: only active workspace)")
    .action((opts: { all?: boolean }) => {
      if (opts.all) {
        clearAll();
        printJson({ ok: true, cleared: "all" });
        return;
      }
      const defaultWs = getDefaultWorkspace();
      if (defaultWs) {
        try {
          removeWorkspace(defaultWs);
        } catch {
          // workspace already gone
        }
      }
      clearApiKey();
      const newDefault = getDefaultWorkspace();
      printJson({
        ok: true,
        removed: defaultWs ?? null,
        remaining_workspaces: Object.keys(getWorkspaces()),
        default_workspace: newDefault ?? null,
      });
    });

  const workspace = auth.command("workspace").description("Manage workspace profiles");

  workspace
    .command("list")
    .description("List all stored workspaces")
    .action(() => {
      const workspaces = getWorkspaces();
      const defaultWs = getDefaultWorkspace();
      const items = Object.entries(workspaces).map(([alias, ws]) => ({
        alias,
        name: ws.name,
        urlKey: ws.urlKey,
        default: alias === defaultWs,
      }));
      printJson({ items });
    });

  workspace
    .command("switch")
    .description("Set default workspace")
    .argument("<alias>", "Workspace alias to switch to")
    .action((alias: string) => {
      try {
        setDefaultWorkspace(alias);
        printJson({ ok: true, default_workspace: alias });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Switch failed");
      }
    });

  workspace
    .command("remove")
    .description("Remove a stored workspace")
    .argument("<alias>", "Workspace alias to remove")
    .action((alias: string) => {
      try {
        const wasDefault = getDefaultWorkspace() === alias;
        removeWorkspace(alias);
        const newDefault = getDefaultWorkspace();
        printJson({
          ok: true,
          removed: alias,
          ...(wasDefault && { warning: "Removed the default workspace" }),
          default_workspace: newDefault ?? null,
        });
      } catch (err) {
        printError(err instanceof Error ? err.message : "Remove failed");
      }
    });
}
