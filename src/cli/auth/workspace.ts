import type { Command } from "commander";
import {
  getWorkspaces,
  getDefaultWorkspace,
  setDefaultWorkspace,
  removeWorkspace,
} from "../../lib/config.ts";
import { printError, printJson } from "../../lib/output.ts";

export function registerWorkspace(auth: Command): void {
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
