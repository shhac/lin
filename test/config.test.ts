import { afterAll, beforeAll, describe, expect, test } from "bun:test";
import { existsSync, mkdirSync, readFileSync, rmSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";

const TEST_CONFIG_DIR = join(tmpdir(), `lin-config-test-${Date.now()}`);

// Set XDG before importing config module so it uses our temp dir
process.env.XDG_CONFIG_HOME = TEST_CONFIG_DIR;
delete process.env.LINEAR_API_KEY;

const { storeApiKey, getApiKey, clearApiKey } = await import("../src/lib/config.ts");

describe("config", () => {
  beforeAll(() => {
    mkdirSync(TEST_CONFIG_DIR, { recursive: true });
  });

  afterAll(() => {
    rmSync(TEST_CONFIG_DIR, { recursive: true, force: true });
    delete process.env.XDG_CONFIG_HOME;
  });

  test("getApiKey returns undefined when no key stored", () => {
    expect(getApiKey()).toBeUndefined();
  });

  test("storeApiKey writes key and getApiKey reads it back", () => {
    storeApiKey("lin_api_test_key_123");
    expect(getApiKey()).toBe("lin_api_test_key_123");
  });

  test("storeApiKey creates config.json in config dir", () => {
    const configPath = join(TEST_CONFIG_DIR, "lin", "config.json");
    expect(existsSync(configPath)).toBe(true);
    const raw = JSON.parse(readFileSync(configPath, "utf8"));
    expect(raw.api_key).toBe("lin_api_test_key_123");
  });

  test("clearApiKey removes the stored key", () => {
    clearApiKey();
    expect(getApiKey()).toBeUndefined();
  });

  test("LINEAR_API_KEY env var takes precedence over stored key", () => {
    storeApiKey("stored_key");
    process.env.LINEAR_API_KEY = "env_key";
    expect(getApiKey()).toBe("env_key");
    delete process.env.LINEAR_API_KEY;
    expect(getApiKey()).toBe("stored_key");
    clearApiKey();
  });
});
