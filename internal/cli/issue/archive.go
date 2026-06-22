package issue

import (
	"context"

	"github.com/spf13/cobra"

	libcli "github.com/shhac/lib-agent-cli/cli"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerArchive(parent *cobra.Command) {
	var archiveYes bool
	archive := &cobra.Command{
		Use:   "archive <id>",
		Short: "Archive an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := libcli.RequireConfirm(archiveYes, "archive issue "+args[0]); err != nil {
				output.WriteError(err)
			}
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueArchive(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"archived": resp.IssueArchive.Success})
		},
	}
	libcli.AddConfirmFlag(archive, &archiveYes)
	parent.AddCommand(archive)

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

	var deleteYes bool
	del := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an issue (move to trash)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := libcli.RequireConfirm(deleteYes, "delete issue "+args[0]+" (move to trash)"); err != nil {
				output.WriteError(err)
			}
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueDelete(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.IssueDelete.Success})
		},
	}
	libcli.AddConfirmFlag(del, &deleteYes)
	parent.AddCommand(del)
}
