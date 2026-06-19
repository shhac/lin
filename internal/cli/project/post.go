package project

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

func registerPost(parent *cobra.Command) {
	post := &cobra.Command{
		Use:   "post",
		Short: "Project updates (health/status posts)",
	}
	parent.AddCommand(post)

	registerPostNew(post)
	registerPostList(post)
	registerPostGet(post)

	output.HandleUnknownCommand(post, "Run 'lin project usage' for available post subcommands")
}
