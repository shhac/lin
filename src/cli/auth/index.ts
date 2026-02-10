import type { Command } from "commander";
import { registerLogin } from "./login.ts";
import { registerLogout } from "./logout.ts";
import { registerStatus } from "./status.ts";
import { registerWorkspace } from "./workspace.ts";

export function registerAuthCommand({ program }: { program: Command }): void {
  const auth = program.command("auth").description("Authentication management");
  registerLogin(auth);
  registerStatus(auth);
  registerLogout(auth);
  registerWorkspace(auth);
}
