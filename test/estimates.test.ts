import { describe, expect, test } from "bun:test";
import { formatEstimateScale, getValidEstimates } from "../src/lib/estimates.ts";

describe("getValidEstimates", () => {
  test("fibonacci base scale", () => {
    expect(getValidEstimates({ type: "fibonacci", allowZero: false, extended: false })).toEqual([
      1, 2, 3, 5, 8, 13,
    ]);
  });

  test("fibonacci extended scale", () => {
    expect(getValidEstimates({ type: "fibonacci", allowZero: false, extended: true })).toEqual([
      1, 2, 3, 5, 8, 13, 21, 34,
    ]);
  });

  test("fibonacci with allowZero", () => {
    expect(getValidEstimates({ type: "fibonacci", allowZero: true, extended: false })).toEqual([
      0, 1, 2, 3, 5, 8, 13,
    ]);
  });

  test("linear base scale", () => {
    expect(getValidEstimates({ type: "linear", allowZero: false, extended: false })).toEqual([
      1, 2, 3, 4, 5,
    ]);
  });

  test("linear extended scale", () => {
    expect(getValidEstimates({ type: "linear", allowZero: false, extended: true })).toEqual([
      1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
    ]);
  });

  test("exponential base scale", () => {
    expect(getValidEstimates({ type: "exponential", allowZero: false, extended: false })).toEqual([
      1, 2, 4, 8, 16,
    ]);
  });

  test("exponential extended scale", () => {
    expect(getValidEstimates({ type: "exponential", allowZero: false, extended: true })).toEqual([
      1, 2, 4, 8, 16, 32, 64,
    ]);
  });

  test("tShirt base scale", () => {
    expect(getValidEstimates({ type: "tShirt", allowZero: false, extended: false })).toEqual([
      1, 2, 3, 4, 5,
    ]);
  });

  test("tShirt extended scale", () => {
    expect(getValidEstimates({ type: "tShirt", allowZero: false, extended: true })).toEqual([
      1, 2, 3, 4, 5, 6,
    ]);
  });

  test("notUsed returns empty", () => {
    expect(getValidEstimates({ type: "notUsed", allowZero: false, extended: false })).toEqual([]);
  });

  test("unknown type returns empty", () => {
    expect(getValidEstimates({ type: "bogus", allowZero: false, extended: false })).toEqual([]);
  });
});

describe("formatEstimateScale", () => {
  test("numeric scale formats as pipe-separated values", () => {
    expect(formatEstimateScale("fibonacci", [1, 2, 3, 5, 8, 13])).toBe("1 | 2 | 3 | 5 | 8 | 13");
  });

  test("tShirt scale includes labels", () => {
    expect(formatEstimateScale("tShirt", [1, 2, 3, 4, 5])).toBe(
      "1 (XS) | 2 (S) | 3 (M) | 4 (L) | 5 (XL)",
    );
  });

  test("tShirt with zero and extended", () => {
    expect(formatEstimateScale("tShirt", [0, 1, 2, 3, 4, 5, 6])).toBe(
      "0 (None) | 1 (XS) | 2 (S) | 3 (M) | 4 (L) | 5 (XL) | 6 (XXL)",
    );
  });
});
