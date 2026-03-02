import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { uploadFiles } from "../../lib/upload.ts";

export function registerUpload(file: Command): void {
  file
    .command("upload")
    .description("Upload files to Linear")
    .argument("<paths...>", "File paths to upload")
    .action(async (paths: string[]) => {
      try {
        const client = getClient();
        const results = await uploadFiles(client, paths);
        printJson(results);
      } catch (err) {
        printError(err instanceof Error ? err.message : "Upload failed");
      }
    });
}
