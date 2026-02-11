import type { Command } from "commander";
import { registerSearch } from "./search.ts";
import { registerList } from "./list.ts";
import { registerGet } from "./get.ts";
import { registerNew } from "./new.ts";
import { registerUpdate } from "./update.ts";
import { registerComment } from "./comment.ts";
import { registerRelation } from "./relation.ts";
import { registerArchive } from "./archive.ts";
import { registerAttachment } from "./attachment.ts";
import { registerUsage } from "./usage.ts";

export function registerIssueCommand({ program }: { program: Command }): void {
  const issue = program.command("issue").description("Issue operations");
  registerSearch(issue);
  registerList(issue);
  registerGet(issue);
  registerNew(issue);
  registerUpdate(issue);
  registerComment(issue);
  registerRelation(issue);
  registerArchive(issue);
  registerAttachment(issue);
  registerUsage(issue);
}
