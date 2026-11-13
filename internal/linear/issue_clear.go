package linear

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// IssueClearField is an IssueUpdateInput field Linear accepts an explicit
// empty value for. Fields outside this set (title, stateId, ...) reject null
// with "Argument Validation Error".
type IssueClearField string

const (
	ClearProject  IssueClearField = "projectId"
	ClearAssignee IssueClearField = "assigneeId"
	ClearCycle    IssueClearField = "cycleId"
	ClearParent   IssueClearField = "parentId"
	ClearDueDate  IssueClearField = "dueDate"
	ClearEstimate IssueClearField = "estimate"
	ClearLabels   IssueClearField = "labelIds"
)

func (f IssueClearField) emptyValue() any {
	if f == ClearLabels {
		return []string{}
	}
	return nil
}

// IssueClear empties one field of an issue. The generated IssueUpdate can't
// express this: every IssueUpdateInput field is omitempty (see
// scripts/fix-omitempty.go), so a nil pointer is dropped rather than sent as
// null. This reuses the generated operation with a hand-built input map.
func IssueClear(ctx context.Context, client graphql.Client, id string, field IssueClearField) (*IssueUpdateResponse, error) {
	req := &graphql.Request{
		OpName: "IssueUpdate",
		Query:  IssueUpdate_Operation,
		Variables: map[string]any{
			"id":    id,
			"input": map[string]any{string(field): field.emptyValue()},
		},
	}
	data := &IssueUpdateResponse{}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: data})
	return data, err
}
