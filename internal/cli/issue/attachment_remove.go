package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerAttachmentRemove(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "remove <attachment-id>",
		Short: "Remove an attachment (works for any source type)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.AttachmentDelete(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"deleted": resp.AttachmentDelete.Success})
		},
	})
}
