import type { Command } from "commander";
import { registerSearch } from "./search.ts";
import { registerList } from "./list.ts";
import { registerGet } from "./get.ts";
import { registerNew } from "./new.ts";
import { registerUpdate } from "./update.ts";
import { registerUsage } from "./usage.ts";

export function registerProjectCommand({ program }: { program: Command }): void {
  const project = program.command("project").description("Project operations");
  registerSearch(project);
  registerList(project);
  registerGet(project);
  registerNew(project);
  registerUpdate(project);
  registerUsage(project);
}
