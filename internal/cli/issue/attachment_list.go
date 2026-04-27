package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

func registerAttachmentList(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "list <issue-id>",
		Short: "List attachments on an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.IssueAttachments(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Issue.Attachments.Nodes))
			for i, a := range resp.Issue.Attachments.Nodes {
				items[i] = map[string]any{
					"id":         a.Id,
					"title":      a.Title,
					"url":        a.Url,
					"subtitle":   a.Subtitle,
					"sourceType": a.SourceType,
				}
			}

			output.PrintJSON(items)
		},
	})
}
