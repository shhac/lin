import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { printPaginated } from "../src/lib/output.ts";

describe("printPaginated", () => {
  const logged: string[] = [];
  const originalLog = console.log;

  beforeEach(() => {
    logged.length = 0;
    console.log = (...args: unknown[]) => {
      logged.push(args.map(String).join(" "));
    };
  });

  afterEach(() => {
    console.log = originalLog;
  });

  test("outputs items array even when empty", () => {
    printPaginated([], { hasNextPage: false });
    const parsed = JSON.parse(logged[0]!);
    expect(parsed.items).toEqual([]);
    expect(parsed.pagination).toBeUndefined();
  });

  test("includes pagination when hasNextPage is true", () => {
    printPaginated([{ id: "a" }], { hasNextPage: true, endCursor: "cursor-abc" });
    const parsed = JSON.parse(logged[0]!);
    expect(parsed.items).toHaveLength(1);
    expect(parsed.pagination).toEqual({ hasMore: true, nextCursor: "cursor-abc" });
  });

  test("omits pagination when hasNextPage is false", () => {
    printPaginated([{ id: "a" }], { hasNextPage: false });
    const parsed = JSON.parse(logged[0]!);
    expect(parsed.items).toHaveLength(1);
    expect(parsed.pagination).toBeUndefined();
  });

  test("uses null for missing endCursor", () => {
    printPaginated([], { hasNextPage: true });
    const parsed = JSON.parse(logged[0]!);
    expect(parsed.pagination).toEqual({ hasMore: true, nextCursor: null });
  });

  test("prunes empty values from items", () => {
    printPaginated([{ id: "a", name: null, status: "" }], { hasNextPage: false });
    const parsed = JSON.parse(logged[0]!);
    expect(parsed.items[0]).toEqual({ id: "a" });
  });
});
