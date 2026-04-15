import { afterAll, beforeAll, describe, expect, test } from "bun:test";
import { execSync } from "node:child_process";
import { mkdirSync, rmSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";

/**
 * Tests the getClient() error path:  when no API key is available,
 * `getClient()` should print a JSON error to stderr and exit with code 1.
 *
 * We cannot import getClient directly for this test because it calls
 * process.exit(1), which would kill the test runner.  Instead we spawn
 * a child process that imports and calls getClient without any key.
 */
describe("getClient", () => {
  const EMPTY_CONFIG_DIR = join(tmpdir(), `lin-client-test-${Date.now()}`);

  beforeAll(() => {
    mkdirSync(EMPTY_CONFIG_DIR, { recursive: true });
  });

  afterAll(() => {
    rmSync(EMPTY_CONFIG_DIR, { recursive: true, force: true });
  });

  test("exits with error JSON when no API key is configured", () => {
    // Inline script that imports getClient in an environment with no key
    const script = `
      import { getClient } from "./src/lib/client.ts";
      getClient();
    `;

    let stderr = "";
    let exitCode = 0;
    try {
      execSync(`bun -e '${script}'`, {
        encoding: "utf8",
        stdio: ["pipe", "pipe", "pipe"],
        env: {
          ...process.env,
          XDG_CONFIG_HOME: EMPTY_CONFIG_DIR,
          LINEAR_API_KEY: "",
          HOME: EMPTY_CONFIG_DIR,
        },
        cwd: join(import.meta.dir, ".."),
      });
    } catch (err: unknown) {
      const execErr = err as { stderr?: string; status?: number };
      stderr = execErr.stderr ?? "";
      exitCode = execErr.status ?? 1;
    }

    expect(exitCode).toBe(1);
    // The error output should contain the JSON error message
    expect(stderr).toContain("Not authenticated");
    expect(stderr).toContain("lin auth login");
  });

  test("error output is valid JSON", () => {
    const script = `
      import { getClient } from "./src/lib/client.ts";
      getClient();
    `;

    let stderr = "";
    try {
      execSync(`bun -e '${script}'`, {
        encoding: "utf8",
        stdio: ["pipe", "pipe", "pipe"],
        env: {
          ...process.env,
          XDG_CONFIG_HOME: EMPTY_CONFIG_DIR,
          LINEAR_API_KEY: "",
          HOME: EMPTY_CONFIG_DIR,
        },
        cwd: join(import.meta.dir, ".."),
      });
    } catch (err: unknown) {
      const execErr = err as { stderr?: string };
      stderr = execErr.stderr ?? "";
    }

    // Extract the JSON line from stderr (may have other output lines from bun)
    const jsonLine = stderr.split("\n").find((line) => line.trim().startsWith("{"));
    expect(jsonLine).toBeDefined();
    const parsed = JSON.parse(jsonLine!);
    expect(parsed).toHaveProperty("error");
    expect(typeof parsed.error).toBe("string");
  });
});
