package linear

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// The generated update functions can't empty a field: every input field is
// omitempty (see scripts/fix-omitempty.go), so a nil pointer is dropped rather
// than sent as null. These reuse each generated operation with a hand-built
// input map, decoding into the generated response type.

func updateRaw[T any](ctx context.Context, client graphql.Client, opName, query, id string, input map[string]any) (*T, error) {
	data := new(T)
	req := &graphql.Request{
		OpName:    opName,
		Query:     query,
		Variables: map[string]any{"id": id, "input": input},
	}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: data})
	return data, err
}

// IssueClear empties one field of an issue.
func IssueClear(ctx context.Context, client graphql.Client, id string, field IssueClearField) (*IssueUpdateResponse, error) {
	return updateRaw[IssueUpdateResponse](ctx, client, "IssueUpdate", IssueUpdate_Operation, id, map[string]any{string(field): field.emptyValue()})
}

// ProjectClearLead removes a project's lead.
func ProjectClearLead(ctx context.Context, client graphql.Client, id string) (*ProjectUpdateResponse, error) {
	return updateRaw[ProjectUpdateResponse](ctx, client, "ProjectUpdate", ProjectUpdate_Operation, id, map[string]any{"leadId": nil})
}

// InitiativeClearOwner removes an initiative's owner.
func InitiativeClearOwner(ctx context.Context, client graphql.Client, id string) (*InitiativeUpdateResponse, error) {
	return updateRaw[InitiativeUpdateResponse](ctx, client, "InitiativeUpdate", InitiativeUpdate_Operation, id, map[string]any{"ownerId": nil})
}
