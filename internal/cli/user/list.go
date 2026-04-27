package user

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/resolvers"
)

func registerList(user *cobra.Command) {
	var teamFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Args:  cobra.NoArgs,
	}
	page := output.AddPageFlags(cmd)

	cmd.Run = func(_ *cobra.Command, _ []string) {
		client := linear.GetClient()
		ctx := context.Background()

		if teamFlag != "" {
			resolved, err := resolvers.ResolveTeam(client, teamFlag)
			if err != nil {
				output.PrintError(err.Error())
			}

			resp, err := linear.TeamMembers(ctx, client, resolved.ID, page.Size(), page.Cursor())
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.Team.Members.Nodes))
			for i, m := range resp.Team.Members.Nodes {
				items[i] = mapTeamMember(m)
			}

			output.PrintPage(items, resp.Team.Members.PageInfo.HasNextPage, resp.Team.Members.PageInfo.EndCursor)
			return
		}

		resp, err := linear.UserList(ctx, client, nil, page.Size(), page.Cursor())
		if err != nil {
			output.HandleGraphQLError(err)
		}

		items := make([]map[string]any, len(resp.Users.Nodes))
		for i, u := range resp.Users.Nodes {
			items[i] = mapUserSummary(u)
		}

		output.PrintPage(items, resp.Users.PageInfo.HasNextPage, resp.Users.PageInfo.EndCursor)
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Filter by team")
	user.AddCommand(cmd)
}

func mapUserSummary(u linear.UserListUsersUserConnectionNodesUser) map[string]any {
	return map[string]any{
		"id":          u.Id,
		"name":        u.Name,
		"email":       u.Email,
		"displayName": u.DisplayName,
	}
}

func mapTeamMember(m linear.TeamMembersTeamMembersUserConnectionNodesUser) map[string]any {
	return map[string]any{
		"id":    m.Id,
		"name":  m.Name,
		"email": m.Email,
	}
}

