import type { Command } from "commander";
import { registerList } from "./list.ts";
import { registerGet } from "./get.ts";

export function registerCycleCommand({ program }: { program: Command }): void {
  const cycle = program.command("cycle").description("Cycle operations");
  registerList(cycle);
  registerGet(cycle);
}
