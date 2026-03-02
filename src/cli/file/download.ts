import { Option, type Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { getApiKey } from "../../lib/config.ts";
import { downloadFile, getOrgId, parseFileUrl } from "../../lib/download.ts";
import { printError, printJson } from "../../lib/output.ts";

type DownloadCmdOpts = {
  output?: string;
  outputDir?: string;
  stdout?: boolean;
  force?: boolean;
};

export function registerDownload(file: Command): void {
  file
    .command("download")
    .description("Download a file from Linear")
    .argument("<url-or-path>", "File URL or path segments (e.g. full URL, org/file UUIDs)")
    .addOption(
      new Option("--output <path>", "Save to specific file path").conflicts(["outputDir", "stdout"]),
    )
    .addOption(
      new Option("--output-dir <dir>", "Save to directory (default: current directory)").conflicts([
        "output",
        "stdout",
      ]),
    )
    .addOption(new Option("--stdout", "Write file content to stdout").conflicts(["output", "outputDir"]))
    .option("--force", "Overwrite existing files")
    .action(async (urlOrPath: string, opts: DownloadCmdOpts) => {
      try {
        const client = getClient();
        const orgId = await getOrgId(client);
        const parsed = parseFileUrl(urlOrPath, orgId);
        const apiKey = getApiKey();
        if (!apiKey) {
          printError("Not authenticated. Run: lin auth login");
          return;
        }
        const result = await downloadFile(parsed.url, {
          apiKey,
          output: opts.output,
          outputDir: opts.outputDir,
          stdout: opts.stdout,
          force: opts.force,
        });
        if (opts.stdout) {
          console.error(JSON.stringify(result));
        } else {
          printJson(result);
        }
      } catch (err) {
        printError(err instanceof Error ? err.message : "Download failed");
      }
    });
}
