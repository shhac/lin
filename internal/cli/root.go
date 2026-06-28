package cli

import (
	"strings"

	"github.com/spf13/cobra"

	agentmcp "github.com/shhac/lib-agent-mcp"
	libcli "github.com/shhac/lib-agent-cli/cli"
	"github.com/shhac/lib-agent-cli/xdg"

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
// --format/--timeout/--debug/--color live in the embedded libcli.Globals (bound
// and resolved by NewRoot); the remaining fields are lin's domain flags. lin's
// config-aware setup (config-default format, truncation, width, client) runs in
// NewRoot's ConfigDefaults hook, and `--format pretty` is opted in per command
// via output.AllowPretty — so lin no longer hand-rolls format validation.
type GlobalFlags struct {
	libcli.Globals // Format, TimeoutMS, Debug

	Expand  string
	Full    bool
	Width   int    // --format pretty card width (0 = auto-detect)
	BaseURL string // hidden; overrides the Linear API base URL for tests
}

func newRootCmd(version string) *cobra.Command {
	g := &GlobalFlags{}

	var root *cobra.Command
	root = libcli.NewRoot(libcli.Options{
		Use:           "lin",
		Short:         "Linear CLI for humans and LLMs",
		Version:       version,
		Globals:       &g.Globals,
		DefaultFormat: output.FormatNDJSON,
		UnknownHint:   "run 'lin usage' for full documentation",
		// lin's config-aware per-run setup runs in the ConfigDefaults hook, which
		// NewRoot invokes (before --format validation) on every command.
		ConfigDefaults: func() { applyConfigDefaults(root, g) },
	})

	// NewRoot binds the shared flags (--format/--timeout/--debug/--color) via
	// Globals; lin registers only its own domain flags here.
	pf := root.PersistentFlags()
	pf.StringVarP(&g.Expand, "expand", "e", "", "Expand truncated fields (comma-separated: description,body,content)")
	pf.BoolVarP(&g.Full, "full", "F", false, "Show full content for all truncated fields")
	pf.IntVar(&g.Width, "width", 0, "Card width for --format pretty (0 = auto-detect terminal)")
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

	// Opt the agent-facing entity groups into the MCP tool surface: each becomes
	// one coarse tool that dispatches its subcommands (with a "help" verb), so the
	// surface is ~one-tool-per-group instead of one-per-leaf. auth/config/usage
	// are deliberately left out — they aren't agent tasks.
	exposeGroups(root, "issue", "project", "initiative", "document", "file",
		"customer", "team", "user", "label", "cycle", "api")

	// Expose the command tree as an MCP server (added last, so it reflects the
	// complete tree and the Expose annotations above). --color/--expose are
	// output-shaping, irrelevant to a tool call, so hide them from the schemas.
	// Expose the cache dir (downloads land in cache/downloads) as a read-only
	// "cache" root, so an MCP client can read back a file `file download` saved —
	// e.g. fs get cache downloads/diagram.png — without filesystem access of its
	// own and without ever seeing the host path.
	root.AddCommand(agentmcp.Command(root,
		agentmcp.WithHiddenFlags("color", "expose"),
		agentmcp.WithFileRoots(xdg.Root("cache", config.CacheDir())),
	))

	return root
}

// exposeGroups opts the named top-level commands into the MCP tool surface. A
// command not found by name is skipped silently — the list is a curation of
// agent-facing groups, not a registration check.
func exposeGroups(root *cobra.Command, names ...string) {
	want := make(map[string]bool, len(names))
	for _, n := range names {
		want[n] = true
	}
	for _, c := range root.Commands() {
		if want[c.Name()] {
			agentmcp.Expose(c)
		}
	}
}

// applyConfigDefaults runs in NewRoot's ConfigDefaults hook (before --format
// validation). It applies lin's persisted config defaults and per-run setup:
// truncation, the config-default format (universal only — pretty stays an
// explicit per-call choice), card width, and the request timeout/client.
func applyConfigDefaults(root *cobra.Command, g *GlobalFlags) {
	cfg := config.Read()

	var maxLen int
	if cfg.Settings != nil && cfg.Settings.Truncation != nil && cfg.Settings.Truncation.MaxLength != nil {
		maxLen = *cfg.Settings.Truncation.MaxLength
	}
	truncation.Configure(truncation.ConfigOpts{Expand: g.Expand, Full: g.Full, MaxLength: maxLen})

	// Apply the persisted default format only when --format wasn't passed; it is
	// always a universal value (the config setter rejects "pretty"), so NewRoot
	// validates it cleanly. Record the effective value for ResolveFormat.
	if !root.PersistentFlags().Changed("format") &&
		cfg.Settings != nil && cfg.Settings.Output != nil && cfg.Settings.Output.DefaultFormat != "" {
		g.Format = cfg.Settings.Output.DefaultFormat
	}
	output.ConfigureFormat(g.Format)
	output.ConfigureWidth(g.Width)

	timeout := g.TimeoutMS
	if timeout == 0 && cfg.Settings != nil && cfg.Settings.Request != nil && cfg.Settings.Request.TimeoutMS != nil {
		timeout = *cfg.Settings.Request.TimeoutMS
	}
	linear.Configure(linear.Options{BaseURL: g.BaseURL, TimeoutMS: timeout, Debug: g.Debug})
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
