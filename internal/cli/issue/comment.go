package issue

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/upload"
)

func registerComment(parent *cobra.Command) {
	comment := &cobra.Command{
		Use:   "comment",
		Short: "Comment operations",
	}
	parent.AddCommand(comment)

	registerCommentList(comment)
	registerCommentNew(comment)
	registerCommentGet(comment)
	registerCommentEdit(comment)
	registerCommentReplies(comment)
}

func registerCommentList(parent *cobra.Command) {
	var (
		limit  string
		cursor string
	)

	cmd := &cobra.Command{
		Use:   "list <issue-id>",
		Short: "List comments on an issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.IssueComments(ctx, client, args[0], pageSize, afterPtr)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Issue.Comments.Nodes))
			for i, c := range resp.Issue.Comments.Nodes {
				m := map[string]any{
					"id":        c.Id,
					"body":      c.Body,
					"createdAt": c.CreatedAt,
					"updatedAt": c.UpdatedAt,
				}
				if c.User != nil {
					m["user"] = map[string]any{"id": c.User.Id, "name": c.User.Name}
				}
				if c.Parent != nil {
					m["parent"] = map[string]any{"id": c.Parent.Id}
				}
				// childCount not available in the IssueComments query; omit to match Go query shape
				items[i] = m
			}

			pi := resp.Issue.Comments.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}

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

func registerCommentGet(parent *cobra.Command) {
	parent.AddCommand(&cobra.Command{
		Use:   "get <comment-id>",
		Short: "Get a specific comment",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			resp, err := linear.CommentGet(ctx, client, args[0])
			if err != nil {
				output.HandleGraphQLError(err)
			}

			c := resp.Comment
			result := map[string]any{
				"id":        c.Id,
				"body":      c.Body,
				"createdAt": c.CreatedAt,
				"updatedAt": c.UpdatedAt,
			}
			if c.User != nil {
				result["user"] = map[string]any{"id": c.User.Id, "name": c.User.Name}
			}
			if c.Issue != nil {
				result["issue"] = map[string]any{"id": c.Issue.Id, "identifier": c.Issue.Identifier}
			}
			if c.Parent != nil {
				result["parent"] = map[string]any{"id": c.Parent.Id}
			}

			output.PrintJSON(result)
		},
	})
}

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

func registerCommentReplies(parent *cobra.Command) {
	var (
		limit  string
		cursor string
	)

	cmd := &cobra.Command{
		Use:   "replies <comment-id>",
		Short: "List replies to a comment",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := linear.GetClient()
			ctx := context.Background()

			pageSize := output.ResolvePageSize(limit)
			afterPtr := output.ResolveCursor(cursor)

			resp, err := linear.CommentReplies(ctx, client, args[0], pageSize, afterPtr)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]any, len(resp.Comment.Children.Nodes))
			for i, c := range resp.Comment.Children.Nodes {
				m := map[string]any{
					"id":        c.Id,
					"body":      c.Body,
					"createdAt": c.CreatedAt,
					"updatedAt": c.UpdatedAt,
				}
				if c.User != nil {
					m["user"] = map[string]any{"id": c.User.Id, "name": c.User.Name}
				}
				items[i] = m
			}

			pi := resp.Comment.Children.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
	parent.AddCommand(cmd)
}
