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
	registerUsage(project)
	output.HandleUnknownCommand(project, "To view a project: lin project get <id>")
}

func strPtr(s string) *string { return &s }

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func intPtr(i int) *int { return &i }
