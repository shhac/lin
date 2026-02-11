import type { Command } from "commander";
import { registerList } from "./list.ts";
import { registerGet } from "./get.ts";
import { registerUsage } from "./usage.ts";

export function registerRoadmapCommand({ program }: { program: Command }): void {
  const roadmap = program.command("roadmap").description("Roadmap operations");
  registerList(roadmap);
  registerGet(roadmap);
  registerUsage(roadmap);
}
