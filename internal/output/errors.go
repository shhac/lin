package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	out "github.com/shhac/lib-agent-output"
	apierrors "github.com/shhac/lin/internal/errors"
)

// PrintError emits a structured error envelope on stderr and exits. The bare
// string is treated as agent-fixable (validation/resolver/usage errors the
// calling agent can correct), so the JSON always carries fixable_by per the
// documented {error, fixable_by, hint} contract.
func PrintError(msg string) {
	WriteErrorTo(Stderr(), apierrors.New(msg, apierrors.FixableByAgent))
	os.Exit(1)
}

func PrintErrorf(format string, args ...any) {
	PrintError(fmt.Sprintf(format, args...))
}

func WriteErrorTo(w io.Writer, err error) {
	out.WriteError(w, err)
}

// WriteError writes a structured error to stderr and exits.
// If err is an *apierrors.APIError, includes fixable_by and hint.
// Otherwise wraps as fixable_by: "agent".
func WriteError(err error) {
	WriteErrorTo(Stderr(), err)
	os.Exit(1)
}

// HandleGraphQLError classifies a genqlient error and writes it as a
// structured error to stderr. Use this for errors from Linear API calls.
func HandleGraphQLError(err error) {
	WriteError(apierrors.ClassifyGraphQLError(err))
}

// HandleUnknownCommand registers a handler for unknown subcommands on a cobra command.
func HandleUnknownCommand(cmd *cobra.Command, hint string) {
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			var names []string
			for _, sub := range cmd.Commands() {
				if sub.Name() != "usage" && sub.Name() != "help" {
					names = append(names, sub.Name())
				}
			}
			msg := fmt.Sprintf("unknown command: %q, valid commands: %s", args[0], strings.Join(names, ", "))
			apiErr := apierrors.New(msg, apierrors.FixableByAgent)
			if hint != "" {
				apiErr = apiErr.WithHint(hint)
			}
			WriteError(apiErr)
		}
		return cmd.Help()
	}
}
