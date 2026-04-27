package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/upload"
)

func registerCommentEdit(parent *cobra.Command) {
	var files []string

	cmd := &cobra.Command{
		Use:   "edit <comment-id> <body>",
		Short: "Edit a comment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			body := args[1]

			if len(files) > 0 {
				uploaded, err := upload.UploadFiles(client, files)
				if err != nil {
					output.PrintError(err.Error())
				}
				body = body + "\n\n" + upload.FormatFileMarkdown(uploaded)
			}

			resp, err := linear.CommentUpdate(ctx, client, args[0], linear.CommentUpdateInput{Body: &body})
			if err != nil {
				output.HandleGraphQLError(err)
			}

			output.PrintJSON(map[string]any{"updated": resp.CommentUpdate.Success})
		},
	}

	cmd.Flags().StringArrayVar(&files, "file", nil, "Attach file (repeatable)")
	parent.AddCommand(cmd)
}
