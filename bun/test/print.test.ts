import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { printError, printJson } from "../src/lib/output.ts";

describe("printJson", () => {
  let logged: string[];
  const originalLog = console.log;

  beforeEach(() => {
    logged = [];
    console.log = (...args: unknown[]) => {
      logged.push(args.map(String).join(" "));
    };
  });

  afterEach(() => {
    console.log = originalLog;
  });

  test("outputs valid JSON to stdout", () => {
    printJson({ a: 1, b: "hello" });
    expect(logged).toHaveLength(1);
    expect(JSON.parse(logged[0]!)).toEqual({ a: 1, b: "hello" });
  });

  test("pretty-prints with 2-space indentation", () => {
    printJson({ key: "value" });
    expect(logged[0]).toContain("\n");
    expect(logged[0]).toContain('  "key"');
  });

  test("prunes null and empty values before printing", () => {
    printJson({ keep: "yes", remove: null, empty: "" });
    const parsed = JSON.parse(logged[0]!);
    expect(parsed).toEqual({ keep: "yes" });
    expect(parsed).not.toHaveProperty("remove");
    expect(parsed).not.toHaveProperty("empty");
  });

  test("handles arrays", () => {
    printJson([1, 2, 3]);
    expect(JSON.parse(logged[0]!)).toEqual([1, 2, 3]);
  });

  test("handles nested objects with pruning", () => {
    printJson({
      project: {
        id: "abc",
        name: "Test",
        lead: null,
        milestones: [],
      },
    });
    const parsed = JSON.parse(logged[0]!);
    expect(parsed).toEqual({ project: { id: "abc", name: "Test" } });
  });

  test("preserves zero and false values", () => {
    printJson({ count: 0, active: false });
    expect(JSON.parse(logged[0]!)).toEqual({ count: 0, active: false });
  });
});

describe("printError", () => {
  let errored: string[];
  const originalError = console.error;

  beforeEach(() => {
    errored = [];
    console.error = (...args: unknown[]) => {
      errored.push(args.map(String).join(" "));
    };
    // Ensure clean state
    process.exitCode = 0;
  });

  afterEach(() => {
    console.error = originalError;
    // Always reset exitCode to 0 so bun test does not report failure
    process.exitCode = 0;
  });

  test("outputs JSON error object to stderr", () => {
    printError("Something went wrong");
    expect(errored).toHaveLength(1);
    expect(JSON.parse(errored[0]!)).toEqual({ error: "Something went wrong" });
    // Reset immediately after assertion
    process.exitCode = 0;
  });

  test("sets process.exitCode to 1", () => {
    expect(process.exitCode).toBe(0);
    printError("fail");
    expect(process.exitCode).toBe(1);
    // Reset immediately after assertion
    process.exitCode = 0;
  });

  test("includes exact error message in JSON", () => {
    printError("Not authenticated. Run: lin auth login <api-key>");
    const parsed = JSON.parse(errored[0]!);
    expect(parsed.error).toBe("Not authenticated. Run: lin auth login <api-key>");
    // Reset immediately after assertion
    process.exitCode = 0;
  });
});
