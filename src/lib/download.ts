import type { LinearClient } from "@linear/sdk";
import { basename, extname, join, resolve } from "node:path";
import { existsSync } from "node:fs";

const UPLOAD_HOST = "uploads.linear.app";
const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

const MIME_TO_EXT: Record<string, string> = {
  "image/png": ".png",
  "image/jpeg": ".jpg",
  "image/gif": ".gif",
  "image/webp": ".webp",
  "image/svg+xml": ".svg",
  "application/pdf": ".pdf",
  "text/plain": ".txt",
  "text/csv": ".csv",
  "application/json": ".json",
  "application/zip": ".zip",
  "video/mp4": ".mp4",
  "audio/mpeg": ".mp3",
};

// eslint-disable-next-line no-control-regex -- intentional: strip control chars from filenames
const UNSAFE_FILENAME_RE = /[<>:"|?*\x00-\x1f]/g;

// ── Types ────────────────────────────────────────────────────────────

export type ParsedFileUrl = {
  url: string;
  orgId: string;
  segments: string[];
};

export type DownloadOpts = {
  output?: string;
  outputDir?: string;
  stdout?: boolean;
  force?: boolean;
};

export type DownloadResult = {
  filename: string;
  path: string;
  size: number;
  contentType: string;
};

type FetchResult = {
  headers: Headers;
  url: string;
  contentType: string;
};

// ── URL Parsing ──────────────────────────────────────────────────────

export function parseFileUrl(input: string, defaultOrgId?: string): ParsedFileUrl {
  let pathname: string;

  if (input.startsWith("http://")) {
    throw new Error("Refusing http:// URL — only https:// is allowed for file downloads.");
  }

  if (input.startsWith("https://")) {
    const url = new URL(input);
    if (url.hostname !== UPLOAD_HOST) {
      throw new Error(`Invalid host: "${url.hostname}". Only ${UPLOAD_HOST} URLs are supported.`);
    }
    ({ pathname } = url);
  } else if (input.startsWith(`${UPLOAD_HOST}/`)) {
    pathname = input.slice(UPLOAD_HOST.length);
  } else {
    pathname = `/${input}`;
  }

  const segments = pathname.split("/").filter(Boolean);

  if (segments.length === 0 || segments.length > 3) {
    throw new Error(`Cannot parse file URL: "${input}". Expected 1-3 UUID path segments.`);
  }

  for (const seg of segments) {
    if (!UUID_RE.test(seg)) {
      throw new Error(
        `Invalid UUID segment: "${seg}". Expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`,
      );
    }
  }

  let orgId: string;
  let fileSegments: string[];

  if (segments.length === 3) {
    orgId = segments[0]!;
    fileSegments = segments.slice(1);
  } else {
    if (!defaultOrgId) {
      throw new Error(
        "Cannot infer organization ID. Provide a full URL with org segment, or authenticate first.",
      );
    }
    orgId = defaultOrgId;
    fileSegments = segments;
  }

  const allSegments = [orgId, ...fileSegments];
  const url = `https://${UPLOAD_HOST}/${allSegments.join("/")}`;

  return { url, orgId, segments: allSegments };
}

// ── Org ID ───────────────────────────────────────────────────────────

export async function getOrgId(client: LinearClient): Promise<string> {
  const org = await client.organization;
  return org.id;
}

// ── Download ─────────────────────────────────────────────────────────

export async function downloadFile(
  url: string,
  opts: DownloadOpts & { apiKey: string },
): Promise<DownloadResult> {
  const res = await fetch(url, {
    headers: { Authorization: opts.apiKey },
    redirect: "follow",
  });

  if (!res.ok) {
    throw new Error(`Download failed: ${res.status} ${res.statusText}`);
  }

  const buffer = await res.arrayBuffer();
  const contentType = res.headers.get("content-type") ?? "application/octet-stream";
  const filename = inferFilename({ headers: res.headers, url, contentType });

  if (opts.stdout) {
    await Bun.write(Bun.stdout, buffer);
    return { filename, path: "<stdout>", size: buffer.byteLength, contentType };
  }

  const destPath = resolveDestPath(filename, { opts, contentType });

  await Bun.write(destPath, buffer);

  return {
    filename: basename(destPath),
    path: resolve(destPath),
    size: buffer.byteLength,
    contentType,
  };
}

// ── Internal helpers ─────────────────────────────────────────────────

function inferFilename({ headers, url, contentType }: FetchResult): string {
  const disposition = headers.get("content-disposition");
  if (disposition) {
    const parsed = parseContentDispositionFilename(disposition);
    if (parsed) {
      return sanitizeFilename(parsed);
    }
  }

  const { pathname } = new URL(url);
  const lastSegment = pathname.split("/").filter(Boolean).pop();
  if (lastSegment) {
    const ext = mimeToExt(contentType);
    if (ext) {
      return `${lastSegment}${ext}`;
    }
    return lastSegment;
  }

  return "download";
}

function parseContentDispositionFilename(header: string): string | null {
  const rfc5987 = header.match(/filename\*\s*=\s*UTF-8''([^\s;]+)/i);
  if (rfc5987?.[1]) {
    try {
      return decodeURIComponent(rfc5987[1]);
    } catch {
      // fall through
    }
  }

  const quoted = header.match(/filename\s*=\s*"([^"]+)"/i);
  if (quoted?.[1]) {
    return quoted[1];
  }

  const unquoted = header.match(/filename\s*=\s*([^\s;]+)/i);
  if (unquoted?.[1]) {
    return unquoted[1];
  }

  return null;
}

function mimeToExt(mime: string): string | null {
  const base = mime.split(";")[0]!.trim().toLowerCase();
  return MIME_TO_EXT[base] ?? null;
}

function sanitizeFilename(name: string): string {
  let clean = name.replace(/^.*[/\\]/, "");
  clean = clean.replace(UNSAFE_FILENAME_RE, "_");
  if (clean.length > 255) {
    clean = clean.slice(0, 255);
  }
  return clean || "download";
}

function resolveDestPath(
  filename: string,
  { opts, contentType }: { opts: DownloadOpts; contentType: string },
): string {
  let destPath: string;

  if (opts.output) {
    destPath = resolve(opts.output);
    const outputExt = extname(destPath).toLowerCase();
    const expectedExt = mimeToExt(contentType);
    if (expectedExt && outputExt && outputExt !== expectedExt) {
      console.error(
        `Warning: output extension "${outputExt}" does not match Content-Type "${contentType}" (expected "${expectedExt}")`,
      );
    }
  } else if (opts.outputDir) {
    if (!existsSync(opts.outputDir)) {
      throw new Error(`Output directory does not exist: "${opts.outputDir}"`);
    }
    destPath = join(opts.outputDir, filename);
  } else {
    destPath = join(process.cwd(), filename);
  }

  if (!opts.force && existsSync(destPath)) {
    throw new Error(`File already exists: "${destPath}". Use --force to overwrite.`);
  }

  return destPath;
}
