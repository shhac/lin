package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/api"
	"github.com/shhac/lin/internal/cli/auth"
	"github.com/shhac/lin/internal/cli/configcmd"
	"github.com/shhac/lin/internal/cli/cycle"
	"github.com/shhac/lin/internal/cli/document"
	"github.com/shhac/lin/internal/cli/file"
	"github.com/shhac/lin/internal/cli/initiative"
	"github.com/shhac/lin/internal/cli/issue"
	"github.com/shhac/lin/internal/cli/label"
	"github.com/shhac/lin/internal/cli/project"
	"github.com/shhac/lin/internal/cli/team"
	"github.com/shhac/lin/internal/cli/usage"
	"github.com/shhac/lin/internal/cli/user"
	"github.com/shhac/lin/internal/config"
	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/truncation"
)

var (
	flagExpand  string
	flagFull    bool
	flagFormat  string
	flagTimeout int
	flagDebug   bool
	flagBaseURL string
)

func newRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "lin",
		Short:         "Linear CLI for humans and LLMs",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
			if err := output.ConfigureFormat(flagFormat); err != nil {
				return err
			}
			timeout := flagTimeout
			if timeout == 0 && cfg.Settings != nil && cfg.Settings.Request != nil && cfg.Settings.Request.TimeoutMS != nil {
				timeout = *cfg.Settings.Request.TimeoutMS
			}
			linear.Configure(linear.Options{
				BaseURL:   flagBaseURL,
				TimeoutMS: timeout,
				Debug:     flagDebug,
			})
			return nil
		},
	}

	root.PersistentFlags().StringVarP(&flagExpand, "expand", "e", "", "Expand truncated fields (comma-separated: description,body,content)")
	root.PersistentFlags().BoolVarP(&flagFull, "full", "F", false, "Show full content for all truncated fields")
	root.PersistentFlags().StringVarP(&flagFormat, "format", "f", "", "Output format: json, yaml, jsonl")
	root.PersistentFlags().IntVarP(&flagTimeout, "timeout", "t", 0, "Request timeout in milliseconds")
	root.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "Log redacted HTTP request records to stderr")
	root.PersistentFlags().StringVar(&flagBaseURL, "base-url", "", "Linear API base URL override for tests")
	_ = root.PersistentFlags().MarkHidden("base-url")

	api.Register(root)
	auth.Register(root)
	project.Register(root)
	initiative.Register(root)
	document.Register(root)
	file.Register(root)
	issue.Register(root)
	team.Register(root)
	user.Register(root)
	label.Register(root)
	cycle.Register(root)
	configcmd.Register(root)
	usage.Register(root)

	output.HandleUnknownCommand(root, "Run 'lin usage' for full documentation")

	return root
}

func Execute(version string) error {
	err := newRootCmd(version).Execute()
	if err != nil {
		var apiErr *apierrors.APIError
		if apierrors.As(err, &apiErr) {
			output.WriteError(apiErr)
		}
		msg := err.Error()
		if strings.Contains(msg, "unknown command") || strings.Contains(msg, "unknown flag") {
			output.WriteError(apierrors.New(msg, apierrors.FixableByAgent).
				WithHint("run 'lin usage' for full documentation"))
		}
		output.WriteError(apierrors.Wrap(err, apierrors.FixableByAgent))
	}
	return err
}
