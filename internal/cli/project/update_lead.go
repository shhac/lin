package project

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/cli/clearable"
	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/resolvers"
)

var leadArg = clearable.Arg{Parent: "lin project update", Name: "lead", Placeholder: "user", Clearable: true}

func registerUpdateLead(parent *cobra.Command) {
	parent.AddCommand(leadArg.Command("Update or remove project lead", setLead, clearLead))
}

func setLead(ctx context.Context, client graphql.Client, projectInput, userInput string) (map[string]any, error) {
	project, err := resolvers.ResolveProject(client, projectInput)
	if err != nil {
		return nil, err
	}
	user, err := resolvers.ResolveUser(client, userInput)
	if err != nil {
		return nil, err
	}
	resp, err := linear.ProjectUpdate(ctx, client, project.ID, linear.ProjectUpdateInput{LeadId: &user.ID})
	if err != nil {
		return nil, err
	}
	result := map[string]any{"updated": resp.ProjectUpdate.Success, "id": project.ID}
	if p := resp.ProjectUpdate.Project; p != nil && p.Lead != nil {
		result["lead"] = map[string]any{"id": p.Lead.Id, "name": p.Lead.Name}
	}
	return result, nil
}

func clearLead(ctx context.Context, client graphql.Client, projectInput string) (map[string]any, error) {
	project, err := resolvers.ResolveProject(client, projectInput)
	if err != nil {
		return nil, err
	}
	resp, err := linear.ProjectClearLead(ctx, client, project.ID)
	if err != nil {
		return nil, err
	}
	result := map[string]any{"updated": resp.ProjectUpdate.Success, "id": project.ID}
	p := resp.ProjectUpdate.Project
	if p == nil {
		return result, nil
	}
	if p.Lead != nil {
		return nil, leadArg.StillSetErr()
	}
	result["cleared"] = "lead"
	return result, nil
}
