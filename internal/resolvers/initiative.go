package resolvers

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/linear"
)

type ResolvedInitiative struct {
	ID     string
	Name   string
	SlugId string
}

func ResolveInitiative(client graphql.Client, input string) (ResolvedInitiative, error) {
	if err := requireInput("initiative", input); err != nil {
		return ResolvedInitiative{}, err
	}
	resp, err := linear.InitiativeGet(ctx(), client, input)
	if err == nil {
		return ResolvedInitiative{
			ID: resp.Initiative.Id, Name: resp.Initiative.Name, SlugId: resp.Initiative.SlugId,
		}, nil
	}
	initiatives, err := linear.FetchAll(func(first int, after *string) ([]linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative, bool, *string, error) {
		resp, err := linear.InitiativeList(ctx(), client, nil, first, after)
		if err != nil {
			return nil, false, nil, err
		}
		return resp.Initiatives.Nodes, resp.Initiatives.PageInfo.HasNextPage, resp.Initiatives.PageInfo.EndCursor, nil
	})
	if err != nil {
		return ResolvedInitiative{}, err
	}
	matches := matchByNameOrID(input, initiatives,
		func(i linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative) string { return i.SlugId },
		func(i linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative) string { return i.Name })
	switch len(matches) {
	case 1:
		match := matches[0]
		return ResolvedInitiative{ID: match.Id, Name: match.Name, SlugId: match.SlugId}, nil
	case 0:
		return ResolvedInitiative{}, initiativeNotFoundErr(input, initiatives)
	default:
		return ResolvedInitiative{}, ambiguousErr("initiative", "initiatives", input, matches, describeInitiative, "use the slug ID or UUID")
	}
}

func describeInitiative(i linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative) string {
	return fmt.Sprintf("%s (slug: %s)", i.Name, i.SlugId)
}

func initiativeNotFoundErr(input string, initiatives []linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative) error {
	names := make([]string, len(initiatives))
	for i, in := range initiatives {
		names[i] = fmt.Sprintf("%s (%s)", in.Name, in.SlugId)
	}
	hint := "none"
	if len(names) > 0 {
		hint = formatChoices(names)
	}
	return apierrors.Newf(apierrors.FixableByAgent, "initiative not found: %q, known initiatives: %s, provide a UUID, slug ID, or exact name", input, hint)
}
