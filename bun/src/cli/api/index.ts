import type { Command } from "commander";
import { getClient } from "../../lib/client.ts";
import { printError, printJson } from "../../lib/output.ts";
import { registerUsage } from "./usage.ts";

export function registerApiCommand({ program }: { program: Command }): void {
  const api = program.command("api").description("Raw GraphQL query against Linear API");

  api
    .command("query")
    .description("Execute a raw GraphQL query")
    .argument("<graphql>", "GraphQL query string")
    .option("--variables <json>", "JSON-encoded variables object")
    .action(async (graphql: string, opts: { variables?: string }) => {
      try {
        const client = getClient();
        const variables = opts.variables ? JSON.parse(opts.variables) : undefined;
        const response = await client.client.rawRequest(graphql, variables);

        if (response.data) {
          printJson(response.data);
        } else {
          printError(response.error ?? "No data returned");
        }
      } catch (err) {
        if (err instanceof SyntaxError) {
          printError(`Invalid --variables JSON: ${err.message}`);
          return;
        }
        printError(err instanceof Error ? err.message : "Query failed");
      }
    });

  registerUsage(api);
}
