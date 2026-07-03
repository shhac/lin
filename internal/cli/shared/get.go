package shared

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"
	libcli "github.com/shhac/lib-agent-cli/cli"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/output/pretty"
)

// GetEntities runs the family's multi-capable get for the lin domain: it sets
// up one client, then resolves each id through getOne and streams the result
// per the shared get contract (NDJSON by default — one record or
// {"@unresolved":…} per id in input order; item-level misses stay on stdout,
// command-level failures bubble to the single sink).
//
// getOne must return the record for an id, or a classified *errors.APIError
// (a lib-agent-output *Error) so a not-found / bad-input becomes an
// @unresolved record rather than aborting the batch. Plain errors (e.g. from
// resolvers that return fmt.Errorf) are treated as command-level failures —
// wrap them with apierrors.Wrap(err, FixableByAgent) before returning when
// the error is item-scoped.
func GetEntities(args []string, getOne func(client graphql.Client, id string) (any, error)) error {
	client := linear.GetClient()
	format := string(output.ResolveFormat(output.FormatNDJSON))
	return libcli.EntityGet(output.Stdout(), format, args, func(id string) (any, error) {
		return getOne(client, id)
	})
}

// RunGet is the standard get entrypoint for a domain that has a pretty card
// renderer: it dispatches to GetEntitiesPretty when --format pretty is active
// (casting each fetched item to the map[string]any the renderer expects) and to
// GetEntities otherwise. Commands with a bespoke pretty path (e.g. issue's
// --full sections) call the two entrypoints directly instead.
func RunGet(
	args []string,
	getOne func(client graphql.Client, id string) (any, error),
	render func(d map[string]any, opts pretty.Options) string,
) error {
	if output.WantsPretty() {
		return GetEntitiesPretty(args, getOne, func(item any, opts pretty.Options) string {
			return render(item.(map[string]any), opts)
		})
	}
	return GetEntities(args, getOne)
}

// GetEntitiesPretty is the --format pretty counterpart to GetEntities: it
// fetches each id with the same getOne, then renders a human-readable card via
// render, stacking multiple cards with a full-width rule between them. Item-level
// misses (classified *APIError) render as a compact error card so the batch
// continues; command-level failures bubble out as the structured stderr error.
func GetEntitiesPretty(
	args []string,
	getOne func(client graphql.Client, id string) (any, error),
	render func(item any, opts pretty.Options) string,
) error {
	client := linear.GetClient()
	opts := output.PrettyOptions()
	w := output.Stdout()
	for i, id := range args {
		if i > 0 {
			_, _ = fmt.Fprintf(w, "\n%s\n\n", pretty.Separator(opts))
		}
		item, err := getOne(client, id)
		if err != nil {
			var apiErr *apierrors.APIError
			if apierrors.As(err, &apiErr) {
				_, _ = fmt.Fprintln(w, pretty.ErrorCard(id, apiErr.Message, apiErr.Hint, opts))
				continue
			}
			return err
		}
		_, _ = fmt.Fprintln(w, render(item, opts))
	}
	return nil
}
