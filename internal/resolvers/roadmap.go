package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

type ResolvedRoadmap struct {
	ID     string
	Name   string
	SlugId string
}

func ResolveRoadmap(client graphql.Client, input string) (ResolvedRoadmap, error) {
	resp, err := linear.RoadmapGet(ctx(), client, input)
	if err == nil {
		return ResolvedRoadmap{
			ID: resp.Roadmap.Id, Name: resp.Roadmap.Name, SlugId: resp.Roadmap.SlugId,
		}, nil
	}
	roadmaps, err := linear.FetchAll(func(first int, after *string) ([]linear.RoadmapListRoadmapsRoadmapConnectionNodesRoadmap, bool, *string, error) {
		resp, err := linear.RoadmapList(ctx(), client, first, after)
		if err != nil {
			return nil, false, nil, err
		}
		return resp.Roadmaps.Nodes, resp.Roadmaps.PageInfo.HasNextPage, resp.Roadmaps.PageInfo.EndCursor, nil
	})
	if err != nil {
		return ResolvedRoadmap{}, err
	}
	lower := strings.ToLower(input)
	for _, r := range roadmaps {
		if r.SlugId == input || strings.ToLower(r.Name) == lower {
			return ResolvedRoadmap{ID: r.Id, Name: r.Name, SlugId: r.SlugId}, nil
		}
	}
	var names []string
	for _, r := range roadmaps {
		names = append(names, fmt.Sprintf("%s (%s)", r.Name, r.SlugId))
	}
	hint := "none"
	if len(names) > 0 {
		hint = strings.Join(names, ", ")
	}
	return ResolvedRoadmap{}, fmt.Errorf("roadmap not found: %q, known roadmaps: %s, provide a UUID, slug ID, or exact name", input, hint)
}
