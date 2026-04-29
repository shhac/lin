package document

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerHistory(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "history <id>",
		Short: "List content edit history for a document",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			doc, err := resolvers.ResolveDocument(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.DocumentContentHistory(ctx, client, doc.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			history := resp.DocumentContentHistory.History
			items := make([]map[string]any, len(history))
			for i, h := range history {
				items[i] = map[string]any{
					"id":                    h.Id,
					"actorIds":              h.ActorIds,
					"contentDataSnapshotAt": h.ContentDataSnapshotAt,
					"createdAt":             h.CreatedAt,
				}
			}

			output.PrintJSON(map[string]any{"items": items})
		},
	}

	parent.AddCommand(cmd)
}
