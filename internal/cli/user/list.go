package user

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/output"
	"github.com/shhac/lin/internal/ptr"
	"github.com/shhac/lin/internal/resolvers"
)

func registerList(user *cobra.Command) {
	var teamFlag string
	var limit string
	var cursor string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			client := linear.GetClient()
			ctx := context.Background()
			pageSize := output.ResolvePageSize(limit)

			after := output.ResolveCursor(cursor)

			if teamFlag != "" {
				resolved, err := resolvers.ResolveTeam(client, teamFlag)
				if err != nil {
					output.PrintError(err.Error())
				}

				resp, err := linear.TeamMembers(ctx, client, resolved.ID, pageSize, after)
				if err != nil {
					output.HandleGraphQLError(err)
				}

				items := make([]map[string]any, len(resp.Team.Members.Nodes))
				for i, m := range resp.Team.Members.Nodes {
					items[i] = mapTeamMember(m)
				}

				pi := resp.Team.Members.PageInfo
				output.PrintPaginated(items, &output.Pagination{
					HasMore:    pi.HasNextPage,
					NextCursor: ptr.Deref(pi.EndCursor),
				})
				return
			}

			resp, err := linear.UserList(ctx, client, nil, pageSize, after)
			if err != nil {
				output.HandleGraphQLError(err)
			}

			items := make([]map[string]any, len(resp.Users.Nodes))
			for i, u := range resp.Users.Nodes {
				items[i] = mapUserSummary(u)
			}

			pi := resp.Users.PageInfo
			output.PrintPaginated(items, &output.Pagination{
				HasMore:    pi.HasNextPage,
				NextCursor: ptr.Deref(pi.EndCursor),
			})
		},
	}

	cmd.Flags().StringVar(&teamFlag, "team", "", "Filter by team")
	cmd.Flags().StringVar(&limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor for next page")
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

