package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/api"
	"github.com/shhac/lin/internal/cli/auth"
	"github.com/shhac/lin/internal/cli/configcmd"
	"github.com/shhac/lin/internal/cli/cycle"
	"github.com/shhac/lin/internal/cli/document"
	"github.com/shhac/lin/internal/cli/issue"
	"github.com/shhac/lin/internal/cli/label"
	"github.com/shhac/lin/internal/cli/project"
	"github.com/shhac/lin/internal/cli/roadmap"
	"github.com/shhac/lin/internal/cli/team"
	"github.com/shhac/lin/internal/cli/usage"
	"github.com/shhac/lin/internal/cli/user"
	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/truncation"
)

var (
	flagExpand string
	flagFull   bool
)

func newRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "lin",
		Short:         "Linear CLI for humans and LLMs",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfg := config.Read()
			var maxLen int
			if cfg.Settings != nil && cfg.Settings.Truncation != nil && cfg.Settings.Truncation.MaxLength != nil {
				maxLen = *cfg.Settings.Truncation.MaxLength
			}
			truncation.Configure(truncation.ConfigOpts{
				Expand:    flagExpand,
				Full:      flagFull,
				MaxLength: maxLen,
			})
		},
	}

	root.PersistentFlags().StringVar(&flagExpand, "expand", "", "Expand truncated fields (comma-separated: description,body,content)")
	root.PersistentFlags().BoolVar(&flagFull, "full", false, "Show full content for all truncated fields")

	api.Register(root)
	auth.Register(root)
	project.Register(root)
	roadmap.Register(root)
	document.Register(root)
	issue.Register(root)
	team.Register(root)
	user.Register(root)
	label.Register(root)
	cycle.Register(root)
	configcmd.Register(root)
	usage.Register(root)

	return root
}

func Execute(version string) error {
	err := newRootCmd(version).Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return err
}
