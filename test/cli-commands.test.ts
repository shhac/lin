import { describe, expect, test } from "bun:test";
import { execSync } from "node:child_process";

/**
 * Tests that each CLI command is properly registered and shows help text.
 * These tests verify the command tree structure without hitting the API.
 */

function runCli(args: string): string {
  try {
    return execSync(`bun run src/index.ts ${args}`, {
      encoding: "utf8",
      cwd: process.cwd(),
      timeout: 10000,
    });
  } catch (err: unknown) {
    const execErr = err as { stdout?: string; stderr?: string };
    // Commander outputs help to stdout on --help, but sometimes to stderr on error
    return (execErr.stdout ?? "") + (execErr.stderr ?? "");
  }
}

describe("CLI command registration", () => {
  test("root shows help with all top-level commands", () => {
    const output = runCli("--help");
    // All resource commands should be listed
    for (const cmd of ["auth", "project", "issue", "team", "user", "label", "cycle", "usage"]) {
      expect(output).toContain(cmd);
    }
  });

  test("--version prints version number", () => {
    const output = runCli("--version");
    expect(output.trim()).toMatch(/^\d+\.\d+\.\d+/);
  });

  test("auth --help shows login and status subcommands", () => {
    const output = runCli("auth --help");
    expect(output).toContain("login");
    expect(output).toContain("status");
  });

  test("project --help shows search, list, get, update subcommands", () => {
    const output = runCli("project --help");
    expect(output).toContain("search");
    expect(output).toContain("list");
    expect(output).toContain("get");
    expect(output).toContain("update");
  });

  test("issue --help shows search, list, get, new, update, comment subcommands", () => {
    const output = runCli("issue --help");
    expect(output).toContain("search");
    expect(output).toContain("list");
    expect(output).toContain("get");
    expect(output).toContain("new");
    expect(output).toContain("update");
    expect(output).toContain("comment");
  });

  test("team --help shows list and get subcommands", () => {
    const output = runCli("team --help");
    expect(output).toContain("list");
    expect(output).toContain("get");
  });

  test("user --help shows list and me subcommands", () => {
    const output = runCli("user --help");
    expect(output).toContain("list");
    expect(output).toContain("me");
  });

  test("label --help shows list subcommand", () => {
    const output = runCli("label --help");
    expect(output).toContain("list");
  });

  test("cycle --help shows list and get subcommands", () => {
    const output = runCli("cycle --help");
    expect(output).toContain("list");
    expect(output).toContain("get");
  });
});

describe("CLI nested command registration", () => {
  test("project get --help shows overview, details, issues", () => {
    const output = runCli("project get --help");
    expect(output).toContain("overview");
    expect(output).toContain("details");
    expect(output).toContain("issues");
  });

  test("project update --help shows title, status, description, lead", () => {
    const output = runCli("project update --help");
    expect(output).toContain("title");
    expect(output).toContain("status");
    expect(output).toContain("description");
    expect(output).toContain("lead");
  });

  test("issue get --help shows overview and comments", () => {
    const output = runCli("issue get --help");
    expect(output).toContain("overview");
    expect(output).toContain("comments");
  });

  test("issue update --help shows title, status, assignee, priority, project, labels, estimate, description", () => {
    const output = runCli("issue update --help");
    expect(output).toContain("title");
    expect(output).toContain("status");
    expect(output).toContain("assignee");
    expect(output).toContain("priority");
    expect(output).toContain("project");
    expect(output).toContain("labels");
    expect(output).toContain("estimate");
    expect(output).toContain("description");
  });

  test("issue comment --help shows new subcommand", () => {
    const output = runCli("issue comment --help");
    expect(output).toContain("new");
  });
});

describe("CLI option registration", () => {
  test("issue list --help shows filter options", () => {
    const output = runCli("issue list --help");
    expect(output).toContain("--project");
    expect(output).toContain("--team");
    expect(output).toContain("--assignee");
    expect(output).toContain("--status");
    expect(output).toContain("--priority");
    expect(output).toContain("--label");
    expect(output).toContain("--cycle");
    expect(output).toContain("--limit");
  });

  test("issue new --help shows required --team and optional flags", () => {
    const output = runCli("issue new --help");
    expect(output).toContain("--team");
    expect(output).toContain("--project");
    expect(output).toContain("--assignee");
    expect(output).toContain("--priority");
    expect(output).toContain("--description");
  });

  test("user list --help shows --team option", () => {
    const output = runCli("user list --help");
    expect(output).toContain("--team");
  });

  test("label list --help shows --team option", () => {
    const output = runCli("label list --help");
    expect(output).toContain("--team");
  });

  test("cycle list --help shows --team required option", () => {
    const output = runCli("cycle list --help");
    expect(output).toContain("--team");
  });

  test("project list --help shows filter and limit options", () => {
    const output = runCli("project list --help");
    expect(output).toContain("--team");
    expect(output).toContain("--status");
    expect(output).toContain("--limit");
  });
});
