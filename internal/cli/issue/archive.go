package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerArchive(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "archive <id>",
		Short: "Archive an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueArchive(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"archived": resp.IssueArchive.Success})
		},
	})

	parent.AddCommand(&cobra.Command{
		Use:   "unarchive <id>",
		Short: "Unarchive an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueUnarchive(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"unarchived": resp.IssueUnarchive.Success})
		},
	})

	parent.AddCommand(&cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an issue (move to trash)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueDelete(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.IssueDelete.Success})
		},
	})
}
