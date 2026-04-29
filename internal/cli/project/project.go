package project

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/shared"
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
	shared.RegisterUsage(project, "project", usageText)
	output.HandleUnknownCommand(project, "To view a project: lin project get <id>")
}
