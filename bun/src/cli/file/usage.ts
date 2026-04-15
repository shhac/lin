import type { Command } from "commander";

const USAGE_TEXT = `lin file — File operations (upload, download)

UPLOAD:
  file upload <paths...>                  Upload one or more files to Linear

DOWNLOAD:
  file download <url-or-path>             Download a file from Linear
    --output <path>                       Save to specific file path
    --output-dir <dir>                    Save to directory (default: CWD)
    --stdout                              Write file content to stdout
    --force                               Overwrite existing files

URL FORMATS:
  Full URL      https://uploads.linear.app/<org>/<uuid>/<uuid>
  Host-relative uploads.linear.app/<org>/<uuid>/<uuid>
  Path only     <org>/<uuid>/<uuid>
  Short path    <uuid>/<uuid>   (org inferred from auth)
  Single UUID   <uuid>          (org inferred from auth)

OUTPUT:
  upload  → [{ filename, assetUrl, contentType }]
  download → { filename, path, size, contentType }
  download --stdout → binary to stdout, metadata JSON to stderr

NOTES:
  --output, --output-dir, and --stdout are mutually exclusive.
  Without --force, download refuses to overwrite existing files.
`;

export function registerUsage(file: Command): void {
  file
    .command("usage")
    .description("Print detailed file command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
