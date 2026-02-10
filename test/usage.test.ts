import { describe, expect, test } from "bun:test";
import { execSync } from "node:child_process";

describe("usage command", () => {
  test("outputs concise docs under 1000 tokens", () => {
    const output = execSync("bun run src/index.ts usage", { encoding: "utf8" });
    // Rough token estimate: ~4 chars per token for English text
    const estimatedTokens = output.length / 4;
    expect(estimatedTokens).toBeLessThan(1000);
    expect(output).toContain("lin");
    expect(output).toContain("auth login");
    expect(output).toContain("issue search");
    expect(output).toContain("project get overview");
  });

  test("includes all resource commands", () => {
    const output = execSync("bun run src/index.ts usage", { encoding: "utf8" });
    for (const cmd of ["auth", "project", "issue", "team", "user", "label", "cycle"]) {
      expect(output).toContain(cmd);
    }
  });
});
