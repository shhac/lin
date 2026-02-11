import type { Command } from "commander";
import { printJson, printError } from "../../lib/output.ts";
import { getSettings } from "../../lib/config.ts";
import { SETTING_DEFS, VALID_KEYS, getNestedValue } from "./index.ts";

export function registerGet(config: Command): void {
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
}
