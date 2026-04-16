package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/mappers"
	"github.com/shhac/lin/internal/output"
)

func registerGet(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get issue details: title, description, status, assignee, labels, relationships",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()
			id := args[0]

			resp, err := linear.IssueGet(ctx, client, id)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			commentsResp, err := linear.IssueComments(ctx, client, id, 250, nil)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			attachResp, err := linear.IssueAttachments(ctx, client, id)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			result := mappers.MapIssueDetail(
				resp.Issue,
				commentsResp.Issue.Comments.Nodes,
				attachResp.Issue.Attachments.Nodes,
			)
			output.PrintJSON(result)
		},
	}

	parent.AddCommand(cmd)
}
