package initiative

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerArchive(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "archive <id>",
		Short: "Archive an initiative",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeArchive(context.Background(), client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"archived": resp.InitiativeArchive.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerUnarchive(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "unarchive <id>",
		Short: "Unarchive an initiative",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeUnarchive(context.Background(), client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"unarchived": resp.InitiativeUnarchive.Success})
		},
	}
	parent.AddCommand(cmd)
}

func registerDelete(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an initiative",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			client := linear.GetClient()

			resolved, err := resolvers.ResolveInitiative(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.InitiativeDelete(context.Background(), client, resolved.ID)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.InitiativeDelete.Success})
		},
	}
	parent.AddCommand(cmd)
}
