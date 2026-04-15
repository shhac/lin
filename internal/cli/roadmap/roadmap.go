package roadmap

import (
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/output"
)

// Register adds the roadmap command group to the parent command.
func Register(parent *cobra.Command) {
	roadmap := &cobra.Command{
		Use:   "roadmap",
		Short: "Roadmap operations",
	}
	output.HandleUnknownCommand(roadmap, "Run 'lin roadmap usage' for help")

	registerList(roadmap)
	registerGet(roadmap)
	registerProjects(roadmap)
	registerUsage(roadmap)

	parent.AddCommand(roadmap)
}
