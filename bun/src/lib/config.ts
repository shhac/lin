import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";
import {
  keychainGet,
  keychainSet,
  keychainDelete,
  keychainDeleteAll,
  KEYCHAIN_SERVICE,
} from "./keychain.ts";

const KEYCHAIN_PLACEHOLDER = "__KEYCHAIN__";

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

type Workspace = {
  api_key: string;
  name?: string;
  urlKey?: string;
};

export type TruncationSettings = {
  maxLength?: number;
};

export type PaginationSettings = {
  defaultPageSize?: number;
};

export type Settings = {
  truncation?: TruncationSettings;
  pagination?: PaginationSettings;
};

type Config = {
  api_key?: string;
  default_workspace?: string;
  workspaces?: Record<string, Workspace>;
  settings?: Settings;
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
  writeFileSync(join(dir, "config.json"), JSON.stringify(config, null, 2), {
    encoding: "utf8",
    mode: 0o600,
  });
}

type ApiKeySource = "environment" | "keychain" | "config";

function resolveApiKey(): { key: string; source: ApiKeySource } | undefined {
  const envKey = process.env.LINEAR_API_KEY?.trim();
  if (envKey) {
    return { key: envKey, source: "environment" };
  }

  const config = readConfig();
  const alias = config.default_workspace;
  if (alias) {
    const keychainKey = keychainGet(`api_key:${alias}`, KEYCHAIN_SERVICE);
    if (keychainKey) {
      return { key: keychainKey, source: "keychain" };
    }
    const ws = config.workspaces?.[alias];
    if (ws && ws.api_key !== KEYCHAIN_PLACEHOLDER) {
      return { key: ws.api_key, source: "config" };
    }
  }
  const legacyKey = config.api_key;
  if (legacyKey && legacyKey !== KEYCHAIN_PLACEHOLDER) {
    return { key: legacyKey, source: "config" };
  }
  return undefined;
}

export function getApiKey(): string | undefined {
  return resolveApiKey()?.key;
}

export function getApiKeySource(): ApiKeySource | undefined {
  return resolveApiKey()?.source;
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

export function getWorkspaces(): Record<string, Workspace> {
  return readConfig().workspaces ?? {};
}

export function getDefaultWorkspace(): string | undefined {
  return readConfig().default_workspace;
}

export function storeWorkspace(alias: string, workspace: Workspace): void {
  const config = readConfig();
  config.workspaces = config.workspaces ?? {};
  config.workspaces[alias] = workspace;
  if (!config.default_workspace) {
    config.default_workspace = alias;
  }
  writeConfig(config);
}

/**
 * Store a workspace login. Attempts to store the API key in the macOS
 * Keychain; on success the config file holds a placeholder instead of the
 * real secret. Falls back to plaintext config on non-macOS or keychain error.
 * Clears legacy top-level `api_key` if present to avoid stale fallback.
 */
export function storeLogin(alias: string, workspace: Workspace): void {
  const config = readConfig();
  delete config.api_key;
  config.workspaces = config.workspaces ?? {};
  const stored = keychainSet({
    account: `api_key:${alias}`,
    value: workspace.api_key,
    service: KEYCHAIN_SERVICE,
  });
  config.workspaces[alias] = {
    ...workspace,
    api_key: stored ? KEYCHAIN_PLACEHOLDER : workspace.api_key,
  };
  if (!config.default_workspace) {
    config.default_workspace = alias;
  }
  writeConfig(config);
}

export function setDefaultWorkspace(alias: string): void {
  const config = readConfig();
  if (!config.workspaces?.[alias]) {
    throw new Error(
      `Unknown workspace: ${alias}. Valid: ${Object.keys(config.workspaces ?? {}).join(", ") || "(none)"}`,
    );
  }
  config.default_workspace = alias;
  writeConfig(config);
}

export function removeWorkspace(alias: string): void {
  const config = readConfig();
  if (!config.workspaces?.[alias]) {
    throw new Error(
      `Unknown workspace: ${alias}. Valid: ${Object.keys(config.workspaces ?? {}).join(", ") || "(none)"}`,
    );
  }
  keychainDelete(`api_key:${alias}`, KEYCHAIN_SERVICE);
  delete config.workspaces[alias];
  if (config.default_workspace === alias) {
    const remaining = Object.keys(config.workspaces);
    config.default_workspace = remaining.length > 0 ? remaining[0] : undefined;
  }
  if (Object.keys(config.workspaces).length === 0) {
    delete config.workspaces;
  }
  writeConfig(config);
}

export function getSettings(): Settings {
  return readConfig().settings ?? {};
}

export function updateSettings(partial: Settings): void {
  const config = readConfig();
  config.settings = { ...config.settings, ...partial };
  writeConfig(config);
}

export function resetSettings(): void {
  const config = readConfig();
  delete config.settings;
  writeConfig(config);
}

export function clearAll(): void {
  keychainDeleteAll(KEYCHAIN_SERVICE);
  writeConfig({});
}
