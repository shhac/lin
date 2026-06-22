package cli

import (
	"strings"

	"github.com/spf13/cobra"

	libcli "github.com/shhac/lib-agent-cli/cli"

	"github.com/shhac/lin/internal/cli/api"
	"github.com/shhac/lin/internal/cli/auth"
	"github.com/shhac/lin/internal/cli/configcmd"
	"github.com/shhac/lin/internal/cli/customer"
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

// GlobalFlags carries lin's persistent flags. The shared
// --format/--timeout/--debug live in the embedded libcli.Globals; the remaining
// fields are lin's domain flags. lin keeps its own flag wording, lenient
// --format parser, and config-aware PersistentPreRunE (its stderr envelopes and
// help text are mirrored by the skill docs), so it registers the shared flags
// itself rather than letting NewRoot bind and validate them.
type GlobalFlags struct {
	libcli.Globals // Format, TimeoutMS, Debug

	Expand  string
	Full    bool
	Width   int    // --format pretty card width (0 = auto-detect)
	BaseURL string // hidden; overrides the Linear API base URL for tests
}

func newRootCmd(version string) *cobra.Command {
	g := &GlobalFlags{}

	root := libcli.NewRoot(libcli.Options{
		Use:           "lin",
		Short:         "Linear CLI for humans and LLMs",
		Version:       version,
		DefaultFormat: output.FormatNDJSON,
		UnknownHint:   "run 'lin usage' for full documentation",
	})

	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		cfg := config.Read()
		var maxLen int
		if cfg.Settings != nil && cfg.Settings.Truncation != nil && cfg.Settings.Truncation.MaxLength != nil {
			maxLen = *cfg.Settings.Truncation.MaxLength
		}
		truncation.Configure(truncation.ConfigOpts{
			Expand:    g.Expand,
			Full:      g.Full,
			MaxLength: maxLen,
		})
		if err := output.ConfigureFormat(cmd, g.Format); err != nil {
			return err
		}
		if err := output.ConfigureColor(g.Color); err != nil {
			return err
		}
		output.ConfigureWidth(g.Width)
		timeout := g.TimeoutMS
		if timeout == 0 && cfg.Settings != nil && cfg.Settings.Request != nil && cfg.Settings.Request.TimeoutMS != nil {
			timeout = *cfg.Settings.Request.TimeoutMS
		}
		linear.Configure(linear.Options{
			BaseURL:   g.BaseURL,
			TimeoutMS: timeout,
			Debug:     g.Debug,
		})
		return nil
	}

	pf := root.PersistentFlags()
	pf.StringVarP(&g.Expand, "expand", "e", "", "Expand truncated fields (comma-separated: description,body,content)")
	pf.BoolVarP(&g.Full, "full", "F", false, "Show full content for all truncated fields")
	pf.StringVarP(&g.Format, "format", "f", "", "Output format: json, yaml, jsonl")
	pf.IntVar(&g.Width, "width", 0, "Card width for --format pretty (0 = auto-detect terminal)")
	pf.IntVarP(&g.TimeoutMS, "timeout", "t", 0, "Request timeout in milliseconds")
	pf.BoolVarP(&g.Debug, "debug", "d", false, "Log redacted HTTP request records to stderr")
	pf.StringVar(&g.Color, "color", "auto", "Colorize output: auto (color when the stream is a terminal), always, or never")
	_ = root.RegisterFlagCompletionFunc("color", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveNoFileComp
	})
	pf.StringVar(&g.BaseURL, "base-url", "", "Linear API base URL override for tests")
	_ = pf.MarkHidden("base-url")

	api.Register(root)
	auth.Register(root)
	project.Register(root)
	initiative.Register(root)
	document.Register(root)
	file.Register(root)
	issue.Register(root)
	customer.Register(root)
	team.Register(root)
	user.Register(root)
	label.Register(root)
	cycle.Register(root)
	configcmd.Register(root)
	usage.Register(root)

	// lin keeps its own unknown-subcommand handler: its message wording and the
	// usage/help skip-list differ from libcli's default, and the skill docs
	// mirror the exact stderr envelope. This overrides the root.RunE NewRoot set.
	output.HandleUnknownCommand(root, "run 'lin usage' for full documentation")

	return root
}

// Execute builds the root command and runs it, rendering any bubbled error as
// the family's structured JSON on stderr exactly once. Most commands pre-render
// and exit on failure (see output.WriteError / HandleGraphQLError), so the only
// errors that reach here are cobra's own (unknown command, bad flag) and the
// PersistentPreRunE format check — all annotated to fixable_by: agent with the
// usage hint to match lin's long-standing contract.
func Execute(version string) error {
	err := newRootCmd(version).Execute()
	if err == nil {
		return nil
	}
	output.WriteErrorTo(output.Stderr(), annotateError(err))
	return err
}

// annotateError upgrades cobra's bare command/flag errors to the structured
// envelope lin documents: an APIError passes through untouched, an unknown
// command/flag gains the usage hint, and anything else is marked agent-fixable.
func annotateError(err error) error {
	var apiErr *apierrors.APIError
	if apierrors.As(err, &apiErr) {
		return apiErr
	}
	msg := err.Error()
	if strings.Contains(msg, "unknown command") || strings.Contains(msg, "unknown flag") {
		return apierrors.New(msg, apierrors.FixableByAgent).
			WithHint("run 'lin usage' for full documentation")
	}
	return apierrors.Wrap(err, apierrors.FixableByAgent)
}
