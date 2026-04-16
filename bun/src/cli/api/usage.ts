import type { Command } from "commander";

const USAGE_TEXT = `lin api — Raw GraphQL escape hatch

SUBCOMMANDS:
  api query <graphql> [--variables <json>]   Execute a raw GraphQL query

ARGUMENTS:
  <graphql>    GraphQL query string (use single quotes to avoid shell escaping)

OPTIONS:
  --variables <json>   JSON-encoded variables object

EXAMPLES:
  lin api query '{ viewer { id name email } }'
  lin api query '{ issue(id: "ENG-123") { id title createdAt completedAt } }'
  lin api query 'query($id: String!) { issue(id: $id) { id title } }' --variables '{"id":"ENG-123"}'

OUTPUT:
  Raw JSON response from Linear's GraphQL API (data field only, empty fields pruned).
  Errors print to stderr as { "error": "..." }.

WHEN TO USE:
  Use this as a last resort when no structured lin command covers your needs.
  Always prefer structured commands (issue get, team states, etc.) — they handle
  pagination, ID resolution, and output formatting automatically.

GRAPHQL DOCS:
  Linear API reference: https://studio.apollographql.com/public/Linear-API/variant/current/home
`;

export function registerUsage(api: Command): void {
  api
    .command("usage")
    .description("Print detailed api command documentation (LLM-optimized)")
    .action(() => {
      console.log(USAGE_TEXT.trim());
    });
}
