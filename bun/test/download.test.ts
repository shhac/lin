import { describe, expect, test } from "bun:test";
import {
  parseFileUrl,
  parseContentDispositionFilename,
  sanitizeFilename,
} from "../src/lib/download.ts";

const ORG = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee";
const FILE1 = "11111111-2222-3333-4444-555555555555";
const FILE2 = "66666666-7777-8888-9999-aaaaaaaaaaaa";

describe("parseFileUrl", () => {
  test("parses full https URL with 3 segments", () => {
    const result = parseFileUrl(`https://uploads.linear.app/${ORG}/${FILE1}/${FILE2}`);
    expect(result.orgId).toBe(ORG);
    expect(result.segments).toEqual([ORG, FILE1, FILE2]);
    expect(result.url).toBe(`https://uploads.linear.app/${ORG}/${FILE1}/${FILE2}`);
  });

  test("parses 2-segment URL with defaultOrgId", () => {
    const result = parseFileUrl(`https://uploads.linear.app/${FILE1}/${FILE2}`, ORG);
    expect(result.orgId).toBe(ORG);
    expect(result.segments).toEqual([ORG, FILE1, FILE2]);
  });

  test("parses 1-segment bare UUID with defaultOrgId", () => {
    const result = parseFileUrl(FILE1, ORG);
    expect(result.orgId).toBe(ORG);
    expect(result.segments).toEqual([ORG, FILE1]);
  });

  test("parses host-prefixed path (no scheme)", () => {
    const result = parseFileUrl(`uploads.linear.app/${ORG}/${FILE1}/${FILE2}`);
    expect(result.orgId).toBe(ORG);
    expect(result.segments).toEqual([ORG, FILE1, FILE2]);
  });

  test("throws on http:// URL", () => {
    expect(() => parseFileUrl(`http://uploads.linear.app/${ORG}/${FILE1}`)).toThrow("http://");
  });

  test("throws on wrong host", () => {
    expect(() => parseFileUrl(`https://evil.com/${ORG}/${FILE1}`)).toThrow("Invalid host");
  });

  test("throws on too many segments", () => {
    expect(() =>
      parseFileUrl(`https://uploads.linear.app/${ORG}/${FILE1}/${FILE2}/${FILE1}`),
    ).toThrow("1-3 UUID");
  });

  test("throws on non-UUID segment", () => {
    expect(() => parseFileUrl(`https://uploads.linear.app/not-a-uuid/${FILE1}`)).toThrow(
      "Invalid UUID",
    );
  });

  test("throws when no defaultOrgId and fewer than 3 segments", () => {
    expect(() => parseFileUrl(`${FILE1}/${FILE2}`)).toThrow("Cannot infer organization");
  });
});

describe("parseContentDispositionFilename", () => {
  test("parses RFC 5987 filename*", () => {
    expect(parseContentDispositionFilename("attachment; filename*=UTF-8''my%20file.pdf")).toBe(
      "my file.pdf",
    );
  });

  test("parses quoted filename", () => {
    expect(parseContentDispositionFilename('attachment; filename="report.csv"')).toBe("report.csv");
  });

  test("parses unquoted filename", () => {
    expect(parseContentDispositionFilename("attachment; filename=photo.jpg")).toBe("photo.jpg");
  });

  test("returns null when no filename", () => {
    expect(parseContentDispositionFilename("attachment")).toBeNull();
  });

  test("prefers RFC 5987 over quoted", () => {
    expect(
      parseContentDispositionFilename(
        "attachment; filename=\"fallback.txt\"; filename*=UTF-8''preferred.txt",
      ),
    ).toBe("preferred.txt");
  });
});

describe("sanitizeFilename", () => {
  test("returns name unchanged when safe", () => {
    expect(sanitizeFilename("report.pdf")).toBe("report.pdf");
  });

  test("strips path prefix", () => {
    expect(sanitizeFilename("/path/to/file.txt")).toBe("file.txt");
    expect(sanitizeFilename("C:\\Users\\file.txt")).toBe("file.txt");
  });

  test("replaces unsafe characters with underscores", () => {
    expect(sanitizeFilename('file<name>:"test".txt')).toBe("file_name___test_.txt");
  });

  test("truncates at 255 characters", () => {
    const long = `${"a".repeat(300)}.txt`;
    expect(sanitizeFilename(long).length).toBe(255);
  });

  test("returns 'download' for empty result", () => {
    expect(sanitizeFilename("")).toBe("download");
  });
});
