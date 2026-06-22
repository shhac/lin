package issue

import (
	"context"

	"github.com/spf13/cobra"

	libcli "github.com/shhac/lib-agent-cli/cli"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerAttachmentRemove(parent *cobra.Command) {
	var yes bool
	cmd := &cobra.Command{
		Use:   "remove <attachment-id>",
		Short: "Remove an attachment (works for any source type)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := libcli.RequireConfirm(yes, "remove attachment "+args[0]); err != nil {
				output.WriteError(err)
			}
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.AttachmentDelete(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.AttachmentDelete.Success})
		},
	}
	libcli.AddConfirmFlag(cmd, &yes)
	parent.AddCommand(cmd)
}
