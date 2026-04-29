package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

// Register adds the api command group to the parent command.
func Register(parent *cobra.Command) {
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Raw GraphQL query against Linear API",
	}
	output.HandleUnknownCommand(apiCmd, "Run 'lin api usage' for help")

	registerQuery(apiCmd)
	shared.RegisterUsage(apiCmd, "api", apiUsageText)

	parent.AddCommand(apiCmd)
}

func registerQuery(apiCmd *cobra.Command) {
	var variables string

	cmd := &cobra.Command{
		Use:   "query <graphql>",
		Short: "Execute a raw GraphQL query",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			query := args[0]
			client := linear.GetRawClient()

			var vars any
			if variables != "" {
				if err := json.Unmarshal([]byte(variables), &vars); err != nil {
					output.PrintError(fmt.Sprintf("Invalid --variables JSON: %s", err.Error()))
					return
				}
			}

			ctx := context.Background()
			data, err := client.RawQuery(ctx, query, vars)
			if err != nil {
				output.PrintError(err.Error())
				return
			}

			var decoded any
			if err := json.Unmarshal(data, &decoded); err != nil {
				output.PrintError(err.Error())
				return
			}
			output.PrintJSON(decoded)
		},
	}

	cmd.Flags().StringVar(&variables, "variables", "", "JSON-encoded variables object")
	apiCmd.AddCommand(cmd)
}

const apiUsageText = `lin api — Raw GraphQL escape hatch

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
  Linear API reference: https://studio.apollographql.com/public/Linear-API/variant/current/home`
