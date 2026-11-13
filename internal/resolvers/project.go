package resolvers

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/filters"
	"github.com/shhac/lin/internal/linear"
)

type ResolvedProject struct {
	ID     string
	Name   string
	SlugId string
}

func ResolveProject(client graphql.Client, input string) (ResolvedProject, error) {
	if err := requireInput("project", input); err != nil {
		return ResolvedProject{}, err
	}
	resp, err := linear.ProjectGet(ctx(), client, input)
	if err == nil {
		return ResolvedProject{
			ID: resp.Project.Id, Name: resp.Project.Name, SlugId: resp.Project.SlugId,
		}, nil
	}
	filter := filters.BuildProjectFilter(input)
	listResp, err := linear.ProjectList(ctx(), client, filter, nameMatchLimit, nil)
	if err != nil {
		return ResolvedProject{}, err
	}
	match, matches, kind := matchByNameOrID(input, listResp.Projects.Nodes,
		func(p linear.ProjectListProjectsProjectConnectionNodesProject) string { return p.SlugId },
		func(p linear.ProjectListProjectsProjectConnectionNodesProject) string { return p.Name })
	switch kind {
	case matchFound:
		return ResolvedProject{ID: match.Id, Name: match.Name, SlugId: match.SlugId}, nil
	case matchNone:
		return ResolvedProject{}, apierrors.Newf(apierrors.FixableByAgent, "project not found: %q, provide a UUID, slug ID, or exact name", input)
	default:
		parts := make([]string, len(matches))
		for i, p := range matches {
			parts[i] = describeProject(p.ProjectSummaryFields)
		}
		return ResolvedProject{}, ambiguousErr("project", "projects", input, parts, "use the slug ID or UUID")
	}
}

func describeProject(p linear.ProjectSummaryFields) string {
	desc := fmt.Sprintf("%s (slug: %s", p.Name, p.SlugId)
	if p.State != "" {
		desc += ", state: " + p.State
	}
	if p.Lead != nil {
		desc += ", lead: " + p.Lead.Name
	}
	return desc + ")"
}
