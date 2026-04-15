import type { Command } from "commander";
import { registerList } from "./list.ts";
import { registerUsage } from "./usage.ts";

export function registerLabelCommand({ program }: { program: Command }): void {
  const label = program.command("label").description("Label operations");
  registerList(label);
  registerUsage(label);
}
