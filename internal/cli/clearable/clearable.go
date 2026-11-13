// Package clearable builds `<field> <id> <value>` update subcommands, optionally
// with --clear to remove the value instead.
package clearable

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
)

// Arg describes a subcommand's value argument.
type Arg struct {
	Parent      string // parent command path, e.g. "lin project update"
	Name        string // subcommand name, e.g. "lead"
	Placeholder string // value placeholder in usage, e.g. "user"
	Clearable   bool
}

type (
	SetFunc   func(ctx context.Context, client graphql.Client, id, value string) (map[string]any, error)
	ClearFunc func(ctx context.Context, client graphql.Client, id string) (map[string]any, error)
)

func (a Arg) Use() string {
	if a.Clearable {
		return fmt.Sprintf("%s <id> [<%s> | --clear]", a.Name, a.Placeholder)
	}
	return fmt.Sprintf("%s <id> <%s>", a.Name, a.Placeholder)
}

// Command builds the subcommand. Arguments are validated before a client is
// created, so a usage mistake never needs credentials. clear may be nil when
// the Arg isn't Clearable.
func (a Arg) Command(short string, set SetFunc, clear ClearFunc) *cobra.Command {
	clearing := false
	cmd := &cobra.Command{
		Use:   a.Use(),
		Short: short,
		Args:  cobra.RangeArgs(1, 2),
		Run: func(_ *cobra.Command, args []string) {
			ctx := context.Background()
			if clearing {
				if err := a.CheckClearArgs(args); err != nil {
					output.WriteError(err)
				}
				printResult(clear(ctx, linear.GetClient(), args[0]))
				return
			}
			value, err := a.Value(args)
			if err != nil {
				output.WriteError(err)
			}
			printResult(set(ctx, linear.GetClient(), args[0], value))
		},
	}
	if a.Clearable {
		cmd.Flags().BoolVar(&clearing, "clear", false, "Remove the "+a.noun())
	}
	return cmd
}

func printResult(result map[string]any, err error) {
	if err != nil {
		output.HandleGraphQLError(err)
	}
	output.PrintJSON(result)
}

func (a Arg) CheckClearArgs(args []string) error {
	if len(args) == 1 {
		return nil
	}
	return apierrors.Newf(apierrors.FixableByAgent, "--clear takes only the ID; drop the %s value", a.Placeholder).
		WithHint(fmt.Sprintf("%s %s <id> --clear", a.Parent, a.Name))
}

// Value returns the value to set. A blank value for a clearable field is
// refused rather than sent: for a resolver-backed field Linear would match it
// against everything, and for the rest it is never what was meant.
func (a Arg) Value(args []string) (string, error) {
	if len(args) != 2 || (a.Clearable && strings.TrimSpace(args[1]) == "") {
		return "", a.missingValueErr()
	}
	return args[1], nil
}

func (a Arg) missingValueErr() error {
	err := apierrors.Newf(apierrors.FixableByAgent, "missing %s", a.Placeholder)
	if !a.Clearable {
		return err.WithHint(a.Parent + " " + a.Use())
	}
	return err.WithHint(fmt.Sprintf("to set it: %s %s <id> <%s>; to remove it: %s %s <id> --clear", a.Parent, a.Name, a.Placeholder, a.Parent, a.Name))
}

// StillSetErr reports a clear Linear acknowledged without applying.
func (a Arg) StillSetErr() error {
	return apierrors.Newf(apierrors.FixableByHuman, "Linear reported success but did not remove the %s", a.noun())
}

func (a Arg) noun() string { return strings.ReplaceAll(a.Name, "-", " ") }
