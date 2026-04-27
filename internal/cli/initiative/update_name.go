package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateName(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "name <id> <new-name>",
		Short: "Update initiative name",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeUpdate(ctx, client, resolved.ID, linear.InitiativeUpdateInput{
				Name: &args[1],
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{
				"id":      resolved.ID,
				"name":    args[1],
				"updated": resp.InitiativeUpdate.Success,
			})
		},
	}
	parent.AddCommand(cmd)
}
