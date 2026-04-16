package project

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

func Register(parent *cobra.Command) {
	project := &cobra.Command{
		Use:   "project",
		Short: "Project operations",
	}
	parent.AddCommand(project)

	registerSearch(project)
	registerList(project)
	registerGet(project)
	registerIssues(project)
	registerNew(project)
	registerUpdate(project)
	registerDelete(project)
	registerUnarchive(project)
	registerUsage(project)
	output.HandleUnknownCommand(project, "To view a project: lin project get <id>")
}

