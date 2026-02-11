import type { Command } from "commander";
import { printJson, printError } from "../../lib/output.ts";
import { getSettings, updateSettings } from "../../lib/config.ts";
import type { Settings } from "../../lib/config.ts";
import { SETTING_DEFS, VALID_KEYS } from "./index.ts";

export function registerSet(config: Command): void {
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
}
