import { describe, expect, test } from "bun:test";
import { getPackageVersion } from "../src/lib/version.ts";

describe("getPackageVersion", () => {
  test("returns a valid semver string", () => {
    const version = getPackageVersion();
    expect(version).toMatch(/^\d+\.\d+\.\d+/);
  });

  test("returns 0.3.1 for current package", () => {
    expect(getPackageVersion()).toBe("0.3.1");
  });
});
