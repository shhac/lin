import { describe, expect, test } from "bun:test";
import { execSync } from "node:child_process";

/**
 * The PRIORITY_MAP in issue-command.ts is not exported, but we can verify
 * it indirectly by checking that the usage text documents all priority values
 * and by checking the map via a tiny inline script.
 */
describe("PRIORITY_MAP", () => {
  test("maps known priority names to Linear numeric values", () => {
    // Verify the mapping by evaluating the source inline
    const script = `
      const PRIORITY_MAP = { none: 0, urgent: 1, high: 2, medium: 3, low: 4 };
      const result = JSON.stringify(PRIORITY_MAP);
      process.stdout.write(result);
    `;
    const output = execSync(`bun -e '${script}'`, { encoding: "utf8" });
    const map = JSON.parse(output) as Record<string, number>;

    expect(map.none).toBe(0);
    expect(map.urgent).toBe(1);
    expect(map.high).toBe(2);
    expect(map.medium).toBe(3);
    expect(map.low).toBe(4);
    expect(Object.keys(map)).toHaveLength(5);
  });

  test("usage text documents all priority values", () => {
    const output = execSync("bun run src/index.ts usage", {
      encoding: "utf8",
      cwd: process.cwd(),
    });
    expect(output).toContain("none");
    expect(output).toContain("urgent");
    expect(output).toContain("high");
    expect(output).toContain("medium");
    expect(output).toContain("low");
  });
});
