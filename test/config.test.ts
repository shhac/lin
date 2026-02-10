import { afterAll, beforeAll, describe, expect, test } from "bun:test";
import { existsSync, mkdirSync, readFileSync, rmSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";

const TEST_CONFIG_DIR = join(tmpdir(), `lin-config-test-${Date.now()}`);

// Set XDG before importing config module so it uses our temp dir
process.env.XDG_CONFIG_HOME = TEST_CONFIG_DIR;
delete process.env.LINEAR_API_KEY;

const {
  storeApiKey,
  getApiKey,
  clearApiKey,
  storeWorkspace,
  storeLogin,
  clearAll,
  getWorkspaces,
  getDefaultWorkspace,
  setDefaultWorkspace,
  removeWorkspace,
} = await import("../src/lib/config.ts");

describe("config", () => {
  beforeAll(() => {
    process.env.XDG_CONFIG_HOME = TEST_CONFIG_DIR;
    mkdirSync(join(TEST_CONFIG_DIR, "lin"), { recursive: true });
    clearAll();
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

describe("workspaces", () => {
  beforeAll(() => {
    process.env.XDG_CONFIG_HOME = TEST_CONFIG_DIR;
    mkdirSync(join(TEST_CONFIG_DIR, "lin"), { recursive: true });
    clearAll();
  });

  afterAll(() => {
    rmSync(TEST_CONFIG_DIR, { recursive: true, force: true });
    delete process.env.XDG_CONFIG_HOME;
  });

  test("getWorkspaces returns empty object when none stored", () => {
    expect(getWorkspaces()).toEqual({});
  });

  test("storeWorkspace stores a workspace and sets it as default", () => {
    storeWorkspace("acme", { api_key: "key_acme", name: "Acme Corp", urlKey: "acme" });
    expect(getWorkspaces()).toEqual({
      acme: { api_key: "key_acme", name: "Acme Corp", urlKey: "acme" },
    });
    expect(getDefaultWorkspace()).toBe("acme");
  });

  test("storeWorkspace does not overwrite existing default", () => {
    storeWorkspace("beta", { api_key: "key_beta", name: "Beta Inc", urlKey: "beta" });
    expect(getDefaultWorkspace()).toBe("acme");
    expect(Object.keys(getWorkspaces())).toEqual(["acme", "beta"]);
  });

  test("getApiKey returns default workspace key when workspaces exist", () => {
    clearApiKey();
    expect(getApiKey()).toBe("key_acme");
  });

  test("getApiKey falls back to legacy api_key when no default workspace", () => {
    // Remove all workspaces, set legacy key
    removeWorkspace("acme");
    removeWorkspace("beta");
    storeApiKey("legacy_key");
    expect(getApiKey()).toBe("legacy_key");
  });

  test("setDefaultWorkspace switches the active workspace", () => {
    clearApiKey();
    storeWorkspace("ws1", { api_key: "key_ws1" });
    storeWorkspace("ws2", { api_key: "key_ws2" });
    expect(getApiKey()).toBe("key_ws1");
    setDefaultWorkspace("ws2");
    expect(getDefaultWorkspace()).toBe("ws2");
    expect(getApiKey()).toBe("key_ws2");
  });

  test("setDefaultWorkspace throws for unknown alias", () => {
    expect(() => setDefaultWorkspace("nonexistent")).toThrow("Unknown workspace: nonexistent");
  });

  test("removeWorkspace removes and reassigns default", () => {
    removeWorkspace("ws2");
    expect(getDefaultWorkspace()).toBe("ws1");
    expect(getWorkspaces()).toEqual({ ws1: { api_key: "key_ws1" } });
  });

  test("removeWorkspace clears default when last workspace removed", () => {
    removeWorkspace("ws1");
    expect(getDefaultWorkspace()).toBeUndefined();
    expect(getWorkspaces()).toEqual({});
  });

  test("removeWorkspace throws for unknown alias", () => {
    expect(() => removeWorkspace("nonexistent")).toThrow("Unknown workspace: nonexistent");
  });

  test("LINEAR_API_KEY env var takes precedence over workspace key", () => {
    storeWorkspace("envtest", { api_key: "ws_key" });
    process.env.LINEAR_API_KEY = "env_override";
    expect(getApiKey()).toBe("env_override");
    delete process.env.LINEAR_API_KEY;
    removeWorkspace("envtest");
  });
});

describe("storeLogin (atomic write)", () => {
  beforeAll(() => {
    process.env.XDG_CONFIG_HOME = TEST_CONFIG_DIR;
    mkdirSync(join(TEST_CONFIG_DIR, "lin"), { recursive: true });
    clearAll();
  });

  afterAll(() => {
    rmSync(TEST_CONFIG_DIR, { recursive: true, force: true });
    delete process.env.XDG_CONFIG_HOME;
  });

  test("storeLogin stores workspace without writing legacy api_key", () => {
    storeLogin("test-org", {
      api_key: "test_key_1",
      name: "Test Org",
      urlKey: "test-org",
    });
    expect(getApiKey()).toBe("test_key_1");
    expect(getDefaultWorkspace()).toBe("test-org");
    expect(getWorkspaces()).toEqual({
      "test-org": { api_key: "test_key_1", name: "Test Org", urlKey: "test-org" },
    });
    // Legacy api_key should NOT be set — workspace is the sole source of truth
    const configPath = join(TEST_CONFIG_DIR, "lin", "config.json");
    const raw = JSON.parse(readFileSync(configPath, "utf8"));
    expect(raw.api_key).toBeUndefined();
  });

  test("storeLogin clears stale legacy api_key", () => {
    // Simulate a pre-workspace config with legacy api_key
    storeApiKey("stale_legacy_key");
    expect(getApiKey()).toBe("test_key_1"); // workspace takes precedence
    // Login again — should clear the legacy key
    storeLogin("test-org", {
      api_key: "test_key_1",
      name: "Test Org",
      urlKey: "test-org",
    });
    const configPath = join(TEST_CONFIG_DIR, "lin", "config.json");
    const raw = JSON.parse(readFileSync(configPath, "utf8"));
    expect(raw.api_key).toBeUndefined();
  });

  test("storeLogin for second workspace preserves first workspace", () => {
    storeLogin("other-org", {
      api_key: "test_key_2",
      name: "Other Org",
      urlKey: "other-org",
    });
    // Default stays as first workspace
    expect(getDefaultWorkspace()).toBe("test-org");
    // Both workspaces present
    const ws = getWorkspaces();
    expect(ws["test-org"]?.api_key).toBe("test_key_1");
    expect(ws["other-org"]?.api_key).toBe("test_key_2");
    // No legacy api_key in config
    const configPath = join(TEST_CONFIG_DIR, "lin", "config.json");
    const raw = JSON.parse(readFileSync(configPath, "utf8"));
    expect(raw.api_key).toBeUndefined();
    // getApiKey returns default workspace key
    expect(getApiKey()).toBe("test_key_1");
  });

  test("old dual-write pattern: storeApiKey then storeWorkspace (regression)", () => {
    // Documents the old buggy pattern — storeApiKey wrote legacy key, storeWorkspace
    // added workspace. The legacy key would get stale on subsequent logins.
    clearAll();
    storeApiKey("test_key_3");
    storeWorkspace("dual-org", {
      api_key: "test_key_3",
      name: "Dual Org",
      urlKey: "dual-org",
    });
    // Both writes succeed when called sequentially
    expect(getApiKey()).toBe("test_key_3");
    expect(getWorkspaces()["dual-org"]?.api_key).toBe("test_key_3");
    // But the legacy key is now a stale duplicate — this is what storeLogin avoids
    const configPath = join(TEST_CONFIG_DIR, "lin", "config.json");
    const raw = JSON.parse(readFileSync(configPath, "utf8"));
    expect(raw.api_key).toBe("test_key_3");
    expect(raw.workspaces?.["dual-org"]?.api_key).toBe("test_key_3");
  });
});

describe("clearAll (logout --all)", () => {
  beforeAll(() => {
    process.env.XDG_CONFIG_HOME = TEST_CONFIG_DIR;
    mkdirSync(join(TEST_CONFIG_DIR, "lin"), { recursive: true });
    clearAll();
  });

  afterAll(() => {
    rmSync(TEST_CONFIG_DIR, { recursive: true, force: true });
    delete process.env.XDG_CONFIG_HOME;
  });

  test("clearAll removes all workspaces and legacy key", () => {
    storeLogin("org-a", { api_key: "key_a", name: "Org A", urlKey: "org-a" });
    storeLogin("org-b", { api_key: "key_b", name: "Org B", urlKey: "org-b" });
    expect(Object.keys(getWorkspaces())).toHaveLength(2);

    clearAll();
    expect(getWorkspaces()).toEqual({});
    expect(getDefaultWorkspace()).toBeUndefined();
    expect(getApiKey()).toBeUndefined();
  });

  test("clearAll on empty config is a no-op", () => {
    clearAll();
    expect(getApiKey()).toBeUndefined();
    expect(getWorkspaces()).toEqual({});
  });
});
