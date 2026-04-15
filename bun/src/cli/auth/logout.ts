import type { Command } from "commander";
import {
  clearApiKey,
  clearAll,
  getDefaultWorkspace,
  getWorkspaces,
  removeWorkspace,
} from "../../lib/config.ts";
import { printJson } from "../../lib/output.ts";

export function registerLogout(auth: Command): void {
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
}
