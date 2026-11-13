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
	match, matches, kind := matchByNameOrID(input, initiatives,
		func(i linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative) string { return i.SlugId },
		func(i linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative) string { return i.Name })
	switch kind {
	case matchFound:
		return ResolvedInitiative{ID: match.Id, Name: match.Name, SlugId: match.SlugId}, nil
	case matchAmbiguous:
		parts := make([]string, len(matches))
		for i, m := range matches {
			parts[i] = fmt.Sprintf("%s (slug: %s)", m.Name, m.SlugId)
		}
		return ResolvedInitiative{}, ambiguousErr("initiative", "initiatives", input, parts, "use the slug ID or UUID")
	}
	var names []string
	for _, i := range initiatives {
		names = append(names, fmt.Sprintf("%s (%s)", i.Name, i.SlugId))
	}
	hint := "none"
	if len(names) > 0 {
		hint = formatChoices(names)
	}
	return ResolvedInitiative{}, apierrors.Newf(apierrors.FixableByAgent, "initiative not found: %q, known initiatives: %s, provide a UUID, slug ID, or exact name", input, hint)
}
