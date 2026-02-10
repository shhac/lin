import { LinearClient } from "@linear/sdk";
import { getApiKey } from "./config.ts";

let cachedClient: LinearClient | undefined;

export function getClient(): LinearClient {
  if (cachedClient) {
    return cachedClient;
  }

  const apiKey = getApiKey();
  if (!apiKey) {
    console.error(JSON.stringify({ error: "Not authenticated. Run: lin auth login <api-key>" }));
    process.exit(1);
  }

  cachedClient = new LinearClient({ apiKey });
  return cachedClient;
}
