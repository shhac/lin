package label

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerSearch(label *cobra.Command) {
	var (
		typeFlag string
		teamFlag string
	)

	cmd := &cobra.Command{
		Use:   "search <text>",
		Short: "Search labels by name (case- and accent-insensitive substring)",
		Args:  cobra.ExactArgs(1),
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, args []string) {
		if err := validateType(typeFlag); err != nil {
			output.WriteError(err)
		}
		if err := rejectTeamForProject(typeFlag, teamFlag); err != nil {
			output.WriteError(err)
		}

		client := linear.GetClient()
		ctx := context.Background()

		if typeFlag == labelTypeProject {
			runProjectLabelList(ctx, client, page, filters.ProjectLabelFilterOpts{Search: args[0]})
			return
		}

		teamID, err := resolvers.ResolveOptionalTeamID(client, teamFlag)
		if err != nil {
			output.PrintError(err.Error())
		}
		runIssueLabelList(ctx, client, page, filters.IssueLabelFilterOpts{Search: args[0]}, teamID)
	}

	addTypeFlag(cmd, &typeFlag)
	cmd.Flags().StringVar(&teamFlag, "team", "", "Restrict search to a single team (issue labels only; key, name, or UUID)")
	label.AddCommand(cmd)
}
