package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

type ResolvedInitiative struct {
	ID     string
	Name   string
	SlugId string
}

func ResolveInitiative(client graphql.Client, input string) (ResolvedInitiative, error) {
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
	lower := strings.ToLower(input)
	for _, i := range initiatives {
		if i.SlugId == input || strings.ToLower(i.Name) == lower {
			return ResolvedInitiative{ID: i.Id, Name: i.Name, SlugId: i.SlugId}, nil
		}
	}
	var names []string
	for _, i := range initiatives {
		names = append(names, fmt.Sprintf("%s (%s)", i.Name, i.SlugId))
	}
	hint := "none"
	if len(names) > 0 {
		hint = strings.Join(names, ", ")
	}
	return ResolvedInitiative{}, fmt.Errorf("initiative not found: %q, known initiatives: %s, provide a UUID, slug ID, or exact name", input, hint)
}
