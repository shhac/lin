import type { LinearClient } from "@linear/sdk";
import { basename } from "node:path";

type UploadedFile = { filename: string; assetUrl: string; contentType: string };

export async function uploadFiles(client: LinearClient, paths: string[]): Promise<UploadedFile[]> {
  const results: UploadedFile[] = [];
  for (const filePath of paths) {
    const file = Bun.file(filePath);
    if (!(await file.exists())) {
      throw new Error(`File not found: ${basename(filePath)}`);
    }
    const filename = basename(filePath);
    const { uploadFile } = await client.fileUpload(file.type, filename, file.size);
    if (!uploadFile) {
      throw new Error(`Upload failed for ${filename}: no upload URL returned`);
    }
    const headers: Record<string, string> = {};
    for (const h of uploadFile.headers) {
      headers[h.key] = h.value;
    }
    const res = await fetch(uploadFile.uploadUrl, {
      method: "PUT",
      headers,
      body: file,
    });
    if (!res.ok) {
      throw new Error(`Upload failed for ${filename}: ${res.status} ${res.statusText}`);
    }
    results.push({ filename, assetUrl: uploadFile.assetUrl, contentType: file.type });
  }
  return results;
}

export function formatFileMarkdown(files: UploadedFile[]): string {
  return files
    .map((f) => {
      const isImage = f.contentType.startsWith("image/");
      return isImage ? `![${f.filename}](${f.assetUrl})` : `[${f.filename}](${f.assetUrl})`;
    })
    .join("\n");
}
