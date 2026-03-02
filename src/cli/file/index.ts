import type { Command } from "commander";
import { handleUnknownCommand } from "../../lib/output.ts";
import { registerDownload } from "./download.ts";
import { registerUpload } from "./upload.ts";
import { registerUsage } from "./usage.ts";

export function registerFileCommand({ program }: { program: Command }): void {
  const file = program.command("file").description("File operations");
  registerDownload(file);
  registerUpload(file);
  registerUsage(file);
  handleUnknownCommand(file, "Upload files: lin file upload <paths...>");
}
