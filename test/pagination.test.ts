import { afterAll, beforeAll, describe, expect, test } from "bun:test";
import { mkdirSync, rmSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";

const TEST_CONFIG_DIR = join(tmpdir(), `lin-pagination-test-${Date.now()}`);

// Set XDG before importing so resolvePageSize reads from our temp dir
process.env.XDG_CONFIG_HOME = TEST_CONFIG_DIR;

const { resolvePageSize } = await import("../src/lib/output.ts");
const { clearAll, updateSettings } = await import("../src/lib/config.ts");

describe("resolvePageSize", () => {
  const configDir = join(TEST_CONFIG_DIR, "lin");

  beforeAll(() => {
    mkdirSync(configDir, { recursive: true });
    clearAll();
  });

  afterAll(() => {
    rmSync(TEST_CONFIG_DIR, { recursive: true, force: true });
    delete process.env.XDG_CONFIG_HOME;
  });

  test("returns parsed --limit when provided", () => {
    expect(resolvePageSize({ limit: "25" })).toBe(25);
  });

  test("returns DEFAULT_PAGE_SIZE (50) when no --limit and no config", () => {
    clearAll();
    expect(resolvePageSize({})).toBe(50);
  });

  test("returns config value when no --limit but config is set", () => {
    updateSettings({ pagination: { defaultPageSize: 100 } });
    expect(resolvePageSize({})).toBe(100);
  });

  test("--limit overrides config value", () => {
    updateSettings({ pagination: { defaultPageSize: 100 } });
    expect(resolvePageSize({ limit: "10" })).toBe(10);
  });
});
