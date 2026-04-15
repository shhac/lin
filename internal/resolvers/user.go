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
	resp, err := linear.UserList(ctx(), client, nil, 250, nil)
	if err != nil {
		return ResolvedUser{}, err
	}
	lower := strings.ToLower(input)
	var matches []ResolvedUser
	for _, u := range resp.Users.Nodes {
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
		for _, u := range resp.Users.Nodes {
			names = append(names, fmt.Sprintf("%s <%s>", u.Name, u.Email))
		}
		return ResolvedUser{}, fmt.Errorf("User not found: %q. Known users: %s", input, strings.Join(names, ", "))
	}
	var ambiguous []string
	for _, u := range matches {
		ambiguous = append(ambiguous, fmt.Sprintf("%s <%s> (%s)", u.Name, u.Email, u.ID))
	}
	return ResolvedUser{}, fmt.Errorf("Ambiguous user: %q matches %d users: %s. Use a unique name, email, or ID.", input, len(matches), strings.Join(ambiguous, ", "))
}
