package resolvers

import (
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
)

type ResolvedWorkflowState struct {
	ID   string
	Name string
}

func ResolveWorkflowState(client graphql.Client, name, teamID string) (ResolvedWorkflowState, error) {
	filter := &linear.WorkflowStateFilter{
		Team: &linear.TeamFilter{Id: &linear.IDComparator{Eq: ptr.To(teamID)}},
	}
	resp, err := linear.WorkflowStates(ctx(), client, filter)
	if err != nil {
		return ResolvedWorkflowState{}, err
	}
	lower := strings.ToLower(name)
	for _, s := range resp.WorkflowStates.Nodes {
		if strings.ToLower(s.Name) == lower {
			return ResolvedWorkflowState{ID: s.Id, Name: s.Name}, nil
		}
	}
	seen := map[string]bool{}
	var validNames []string
	for _, s := range resp.WorkflowStates.Nodes {
		if !seen[s.Name] {
			seen[s.Name] = true
			validNames = append(validNames, s.Name)
		}
	}
	return ResolvedWorkflowState{}, fmt.Errorf("Unknown status: %q. Valid values: %s", name, strings.Join(validNames, " | "))
}
