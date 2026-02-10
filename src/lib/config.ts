import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";

function getConfigDir(): string {
  const xdg = process.env.XDG_CONFIG_HOME?.trim();
  if (xdg) {
    return join(xdg, "lin");
  }
  return join(homedir(), ".config", "lin");
}

function ensureConfigDir(): string {
  const dir = getConfigDir();
  if (!existsSync(dir)) {
    mkdirSync(dir, { recursive: true });
  }
  return dir;
}

type Config = {
  api_key?: string;
};

function readConfig(): Config {
  const configPath = join(getConfigDir(), "config.json");
  if (!existsSync(configPath)) {
    return {};
  }
  try {
    return JSON.parse(readFileSync(configPath, "utf8")) as Config;
  } catch {
    return {};
  }
}

function writeConfig(config: Config): void {
  const dir = ensureConfigDir();
  writeFileSync(join(dir, "config.json"), JSON.stringify(config, null, 2), "utf8");
}

export function getApiKey(): string | undefined {
  const envKey = process.env.LINEAR_API_KEY?.trim();
  if (envKey) {
    return envKey;
  }
  return readConfig().api_key;
}

export function storeApiKey(key: string): void {
  const config = readConfig();
  config.api_key = key;
  writeConfig(config);
}

export function clearApiKey(): void {
  const config = readConfig();
  delete config.api_key;
  writeConfig(config);
}
