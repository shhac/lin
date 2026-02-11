import type { Command } from "commander";
import { printJson, printError } from "../lib/output.ts";
import { getSettings, updateSettings, resetSettings } from "../lib/config.ts";
import type { Settings } from "../lib/config.ts";

type SettingDef = {
  path: [keyof Settings, string];
  parse: (v: string) => unknown;
  description: string;
};

const SETTING_DEFS: Record<string, SettingDef> = {
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

const VALID_KEYS = Object.keys(SETTING_DEFS);

function getNestedValue(settings: Settings, key: string): unknown {
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

  config
    .command("get")
    .argument("[key]", "Setting key (omit to show all)")
    .description("Show current settings")
    .action((key?: string) => {
      const settings = getSettings();

      if (!key) {
        printJson(settings);
        return;
      }

      if (!SETTING_DEFS[key]) {
        printError(`Unknown setting: ${key}. Valid keys: ${VALID_KEYS.join(", ")}`);
        return;
      }

      const value = getNestedValue(settings, key);
      printJson({ [key]: value ?? null });
    });

  config
    .command("set")
    .argument("<key>", "Setting key")
    .argument("<value>", "Setting value")
    .description("Update a setting")
    .action((key: string, value: string) => {
      const def = SETTING_DEFS[key];
      if (!def) {
        printError(`Unknown setting: ${key}. Valid keys: ${VALID_KEYS.join(", ")}`);
        return;
      }

      let parsed: unknown;
      try {
        parsed = def.parse(value);
      } catch (err) {
        printError((err as Error).message);
        return;
      }

      const [section, field] = def.path;
      const settings = getSettings();
      const sectionObj = (settings[section] ?? {}) as Record<string, unknown>;
      sectionObj[field] = parsed;
      updateSettings({ [section]: sectionObj } as Settings);
      printJson({ [key]: parsed });
    });

  config
    .command("reset")
    .argument("[key]", "Setting key (omit to reset all)")
    .description("Reset settings to defaults")
    .action((key?: string) => {
      if (!key) {
        resetSettings();
        printJson({ reset: "all" });
        return;
      }

      const def = SETTING_DEFS[key];
      if (!def) {
        printError(`Unknown setting: ${key}. Valid keys: ${VALID_KEYS.join(", ")}`);
        return;
      }

      const [section, field] = def.path;
      const settings = getSettings();
      const sectionObj = settings[section] as Record<string, unknown> | undefined;
      if (sectionObj) {
        delete sectionObj[field];
        updateSettings({ [section]: sectionObj } as Settings);
      }
      printJson({ reset: key });
    });

  config
    .command("list-keys")
    .description("List all available setting keys")
    .action(() => {
      const keys = Object.entries(SETTING_DEFS).map(([key, def]) => ({
        key,
        description: def.description,
      }));
      printJson({ keys });
    });
}
