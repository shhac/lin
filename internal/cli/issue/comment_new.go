package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/upload"
)

func registerCommentNew(parent *cobra.Command) {
	var (
		parentComment string
		files         []string
	)

	cmd := &cobra.Command{
		Use:   "new <issue-id> <body>",
		Short: "Add comment to issue",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			issueID := args[0]
			body := args[1]

			if len(files) > 0 {
				uploaded, err := upload.UploadFiles(client, files)
				if err != nil {
					output.PrintError(err.Error())
				}
				body = body + "\n\n" + upload.FormatFileMarkdown(uploaded)
			}

			input := linear.CommentCreateInput{
				IssueId: &issueID,
				Body:    &body,
			}
			if parentComment != "" {
				input.ParentId = &parentComment
			}

			resp, err := linear.CommentCreate(ctx, client, input)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{
				"id":      resp.CommentCreate.Comment.Id,
				"body":    resp.CommentCreate.Comment.Body,
				"created": resp.CommentCreate.Success,
			})
		},
	}

	cmd.Flags().StringVar(&parentComment, "parent", "", "Parent comment ID (threaded reply)")
	cmd.Flags().StringArrayVar(&files, "file", nil, "Attach file (repeatable)")
	parent.AddCommand(cmd)
}
