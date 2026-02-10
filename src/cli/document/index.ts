import type { Command } from "commander";
import { registerSearch } from "./search.ts";
import { registerList } from "./list.ts";
import { registerGet } from "./get.ts";
import { registerNew } from "./new.ts";
import { registerUpdate } from "./update.ts";

export function registerDocumentCommand({ program }: { program: Command }): void {
  const document = program.command("document").description("Document operations");
  registerSearch(document);
  registerList(document);
  registerGet(document);
  registerNew(document);
  registerUpdate(document);
}
