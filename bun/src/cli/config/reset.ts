import type { Command } from "commander";
import { printJson, printError } from "../../lib/output.ts";
import { getSettings, updateSettings, resetSettings } from "../../lib/config.ts";
import type { Settings } from "../../lib/config.ts";
import { SETTING_DEFS, VALID_KEYS } from "./index.ts";

export function registerReset(config: Command): void {
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
}
