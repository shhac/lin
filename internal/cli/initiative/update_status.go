package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateStatus(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "status <id> <new-status>",
		Short: "Update initiative status",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			normalized, err := validateInitiativeStatus(args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			s := linear.InitiativeStatus(normalized)
			resp, err := linear.InitiativeUpdate(ctx, client, resolved.ID, linear.InitiativeUpdateInput{
				Status: &s,
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.InitiativeUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
