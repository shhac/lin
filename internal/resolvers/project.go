package resolvers

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
)

type ResolvedProject struct {
	ID     string
	Name   string
	SlugId string
}

func ResolveProject(client graphql.Client, input string) (ResolvedProject, error) {
	resp, err := linear.ProjectGet(ctx(), client, input)
	if err == nil {
		return ResolvedProject{
			ID: resp.Project.Id, Name: resp.Project.Name, SlugId: resp.Project.SlugId,
		}, nil
	}
	filter := filters.BuildProjectFilter(input)
	listResp, err := linear.ProjectList(ctx(), client, filter, 1, nil)
	if err != nil {
		return ResolvedProject{}, err
	}
	if len(listResp.Projects.Nodes) == 0 {
		return ResolvedProject{}, fmt.Errorf("Project not found: %q. Provide a UUID, slug ID, or exact name.", input)
	}
	p := listResp.Projects.Nodes[0]
	return ResolvedProject{ID: p.Id, Name: p.Name, SlugId: p.SlugId}, nil
}
