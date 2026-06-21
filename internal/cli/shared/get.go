package shared

import (
	"github.com/Khan/genqlient/graphql"
	libcli "github.com/shhac/lib-agent-cli/cli"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
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
