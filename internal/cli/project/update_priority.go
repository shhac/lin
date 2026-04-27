package project

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerUpdatePriority(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "priority <id> <priority>",
		Short: "Update project priority",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := priorities.Resolve(args[1])
			if !ok {
				output.PrintErrorf("Invalid priority: %q. Valid values: %s", args[1], priorities.Values)
			}

			client := linear.GetClient()
			ctx := context.Background()

			resolved, err := resolvers.ResolveProject(client, args[0])
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.ProjectUpdate(ctx, client, resolved.ID, linear.ProjectUpdateInput{
				Priority: ptr.To(p),
			})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.ProjectUpdate.Success})
		},
	}
	parent.AddCommand(cmd)
}
