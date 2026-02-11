import type { Command } from "commander";
import type { Settings } from "../../lib/config.ts";
import { registerGet } from "./get.ts";
import { registerSet } from "./set.ts";
import { registerReset } from "./reset.ts";
import { registerListKeys } from "./list-keys.ts";

export type SettingDef = {
  path: [keyof Settings, string];
  parse: (v: string) => unknown;
  description: string;
};

export const SETTING_DEFS: Record<string, SettingDef> = {
  "truncation.maxLength": {
    path: ["truncation", "maxLength"],
    parse: (v) => {
      const n = Number(v);
      if (!Number.isInteger(n) || n < 0) {
        throw new Error(`Invalid value: ${v}. Must be a non-negative integer.`);
      }
      return n;
    },
    description: "Max characters before truncating description/body/content fields (default: 200)",
  },
  "pagination.defaultPageSize": {
    path: ["pagination", "defaultPageSize"],
    parse: (v) => {
      const n = Number(v);
      if (!Number.isInteger(n) || n < 1 || n > 250) {
        throw new Error(`Invalid value: ${v}. Must be an integer between 1 and 250.`);
      }
      return n;
    },
    description: "Default number of results for list/search commands (default: 50)",
  },
};

export const VALID_KEYS = Object.keys(SETTING_DEFS);

export function getNestedValue(settings: Settings, key: string): unknown {
  const def = SETTING_DEFS[key];
  if (!def) {
    return undefined;
  }
  const [section, field] = def.path;
  const sectionObj = settings[section] as Record<string, unknown> | undefined;
  return sectionObj?.[field];
}

export function registerConfigCommand({ program }: { program: Command }): void {
  const config = program.command("config").description("View and update CLI settings");
  registerGet(config);
  registerSet(config);
  registerReset(config);
  registerListKeys(config);
}
