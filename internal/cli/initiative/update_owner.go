package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdateOwner(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "owner <id> <user>",
		Short: "Update initiative owner",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			user, err := resolvers.ResolveUser(client, args[1])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeUpdate(ctx, client, resolved.ID, linear.InitiativeUpdateInput{
				OwnerId: &user.ID,
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.InitiativeUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
