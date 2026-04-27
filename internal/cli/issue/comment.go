package issue

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

// commentBase builds the common fields shared by all comment output shapes.
func commentBase(id, body, createdAt, updatedAt string, userID, userName *string) map[string]any {
	m := map[string]any{
		"id":        id,
		"body":      body,
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}
	if userID != nil {
		m["user"] = map[string]any{"id": *userID, "name": *userName}
	}
	return m
}

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

	output.HandleUnknownCommand(comment, "Run 'lin issue usage' for available comment subcommands")
}
