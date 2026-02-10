import { describe, expect, test } from "bun:test";
import { pruneEmpty } from "../src/lib/output.ts";

// pruneEmpty returns the same type T, but at runtime removes empty values.
// Cast inputs to `unknown` so TypeScript doesn't complain about structural mismatches.
const prune = (v: unknown) => pruneEmpty(v);

describe("pruneEmpty", () => {
  test("removes null and undefined values", () => {
    expect(prune({ a: 1, b: null, c: undefined })).toEqual({ a: 1 });
  });

  test("removes empty strings", () => {
    expect(prune({ a: "hello", b: "", c: "  " })).toEqual({ a: "hello" });
  });

  test("preserves zero and false", () => {
    expect(prune({ a: 0, b: false, c: true })).toEqual({ a: 0, b: false, c: true });
  });

  test("removes empty arrays", () => {
    expect(prune({ a: [1, 2], b: [] })).toEqual({ a: [1, 2] });
  });

  test("prunes nested objects recursively", () => {
    expect(prune({ a: { b: null, c: { d: "" } }, e: "keep" })).toEqual({ e: "keep" });
  });

  test("prunes null entries from arrays", () => {
    expect(prune([1, null, "hello", undefined, ""])).toEqual([1, "hello"]);
  });

  test("returns empty object for fully pruned input", () => {
    expect(prune({ a: null, b: "" })).toEqual({});
  });

  test("handles deeply nested structures", () => {
    const input = {
      project: {
        id: "abc",
        name: "Test",
        lead: null,
        milestones: [
          { id: "m1", name: "Alpha", targetDate: null },
          { id: "", name: "", targetDate: null },
        ],
      },
    };
    expect(prune(input)).toEqual({
      project: {
        id: "abc",
        name: "Test",
        milestones: [{ id: "m1", name: "Alpha" }],
      },
    });
  });
});
