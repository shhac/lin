package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
)

type ResolvedUser struct {
	ID          string
	Name        string
	Email       string
	DisplayName string
}

func ResolveUser(client graphql.Client, input string) (ResolvedUser, error) {
	users, err := linear.FetchAll(func(first int, after *string) ([]linear.UserListUsersUserConnectionNodesUser, bool, *string, error) {
		resp, err := linear.UserList(ctx(), client, nil, first, after)
		if err != nil {
			return nil, false, nil, err
		}
		return resp.Users.Nodes, resp.Users.PageInfo.HasNextPage, resp.Users.PageInfo.EndCursor, nil
	})
	if err != nil {
		return ResolvedUser{}, err
	}
	lower := strings.ToLower(input)
	var matches []ResolvedUser
	for _, u := range users {
		if u.Id == input ||
			strings.ToLower(u.Name) == lower ||
			strings.ToLower(u.Email) == lower ||
			strings.ToLower(u.DisplayName) == lower {
			matches = append(matches, ResolvedUser{
				ID: u.Id, Name: u.Name, Email: u.Email, DisplayName: u.DisplayName,
			})
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) == 0 {
		var names []string
		for _, u := range users {
			names = append(names, fmt.Sprintf("%s <%s>", u.Name, u.Email))
		}
		return ResolvedUser{}, fmt.Errorf("user not found: %q, known users: %s", input, strings.Join(names, ", "))
	}
	var ambiguous []string
	for _, u := range matches {
		ambiguous = append(ambiguous, fmt.Sprintf("%s <%s> (%s)", u.Name, u.Email, u.ID))
	}
	return ResolvedUser{}, fmt.Errorf("ambiguous user: %q matches %d users: %s, use a unique name, email, or ID", input, len(matches), strings.Join(ambiguous, ", "))
}
