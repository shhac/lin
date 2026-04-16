package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
)

type ResolvedTeam struct {
	ID   string
	Name string
	Key  string
}

func ResolveTeam(client graphql.Client, input string) (ResolvedTeam, error) {
	resp, err := linear.TeamGet(ctx(), client, input)
	if err == nil {
		return ResolvedTeam{ID: resp.Team.Id, Name: resp.Team.Name, Key: resp.Team.Key}, nil
	}
	filter := filters.BuildTeamFilter(input)
	listResp, err := linear.TeamList(ctx(), client, filter, 50, nil)
	if err != nil {
		return ResolvedTeam{}, err
	}
	if len(listResp.Teams.Nodes) > 0 {
		t := listResp.Teams.Nodes[0]
		return ResolvedTeam{ID: t.Id, Name: t.Name, Key: t.Key}, nil
	}
	allResp, err := linear.TeamList(ctx(), client, nil, 250, nil)
	if err != nil {
		return ResolvedTeam{}, err
	}
	var keys []string
	for _, t := range allResp.Teams.Nodes {
		keys = append(keys, fmt.Sprintf("%s (%s)", t.Key, t.Name))
	}
	hint := "none"
	if len(keys) > 0 {
		hint = strings.Join(keys, ", ")
	}
	return ResolvedTeam{}, fmt.Errorf("team not found: %q, known teams: %s, provide a UUID, key, or exact name", input, hint)
}
