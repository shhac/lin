package issue

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

func Register(parent *cobra.Command) {
	issue := &cobra.Command{
		Use:   "issue",
		Short: "Issue operations",
	}
	parent.AddCommand(issue)

	registerSearch(issue)
	registerList(issue)
	registerGet(issue)
	registerNew(issue)
	registerUpdate(issue)
	registerComment(issue)
	registerRelation(issue)
	registerArchive(issue)
	registerAttachment(issue)
	registerHistory(issue)
	registerUsage(issue)
	output.HandleUnknownCommand(issue, "To view an issue: lin issue get <id>")
}
