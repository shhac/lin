import { describe, expect, test } from "bun:test";
import { execSync } from "node:child_process";

describe("usage command", () => {
  test("outputs concise docs under 1000 tokens", () => {
    const output = execSync("bun run src/index.ts usage", { encoding: "utf8" });
    // Rough token estimate: ~4 chars per token for English text
    const estimatedTokens = output.length / 4;
    expect(estimatedTokens).toBeLessThan(1100);
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

  test("mentions sub-usage discovery pattern", () => {
    const output = execSync("bun run src/index.ts usage", { encoding: "utf8" });
    expect(output).toContain("<command> usage");
  });
});

describe("sub-command usage", () => {
  const subCommands = [
    {
      cmd: "issue",
      mustContain: [
        "issue search",
        "issue list",
        "issue new",
        "issue update",
        "issue comment",
        "issue relation",
        "issue archive",
        "issue attachment",
      ],
    },
    {
      cmd: "document",
      mustContain: [
        "document search",
        "document list",
        "document get",
        "document new",
        "document update",
      ],
    },
    {
      cmd: "project",
      mustContain: [
        "project search",
        "project list",
        "project get",
        "project new",
        "project update",
      ],
    },
    { cmd: "roadmap", mustContain: ["roadmap list", "roadmap get"] },
    { cmd: "team", mustContain: ["team list", "team get", "team states"] },
    { cmd: "user", mustContain: ["user list", "user me"] },
    { cmd: "label", mustContain: ["label list"] },
    { cmd: "cycle", mustContain: ["cycle list", "cycle get"] },
    { cmd: "auth", mustContain: ["auth login", "auth status"] },
    { cmd: "config", mustContain: ["config get", "config set"] },
  ] as const;

  for (const { cmd, mustContain } of subCommands) {
    test(`${cmd} usage is under 500 tokens`, () => {
      const output = execSync(`bun run src/index.ts ${cmd} usage`, { encoding: "utf8" });
      const estimatedTokens = output.length / 4;
      expect(estimatedTokens).toBeLessThan(500);
    });

    test(`${cmd} usage contains key subcommands`, () => {
      const output = execSync(`bun run src/index.ts ${cmd} usage`, { encoding: "utf8" });
      expect(output).toContain(`lin ${cmd}`);
      for (const phrase of mustContain) {
        expect(output).toContain(phrase);
      }
    });
  }
});
