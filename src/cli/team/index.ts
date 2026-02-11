import type { Command } from "commander";
import { registerList } from "./list.ts";
import { registerGet } from "./get.ts";
import { registerStates } from "./states.ts";
import { registerUsage } from "./usage.ts";

export function registerTeamCommand({ program }: { program: Command }): void {
  const team = program.command("team").description("Team operations");
  registerList(team);
  registerGet(team);
  registerStates(team);
  registerUsage(team);
}
