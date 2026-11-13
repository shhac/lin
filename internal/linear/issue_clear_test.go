package linear_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/testutil"
)

func TestIssueClear_SendsExplicitEmptyValue(t *testing.T) {
	cases := []struct {
		field linear.IssueClearField
		want  any
	}{
		{linear.ClearProject, nil},
		{linear.ClearAssignee, nil},
		{linear.ClearCycle, nil},
		{linear.ClearParent, nil},
		{linear.ClearDueDate, nil},
		{linear.ClearEstimate, nil},
		{linear.ClearLabels, []any{}},
	}
	for _, tc := range cases {
		t.Run(string(tc.field), func(t *testing.T) {
			mock := testutil.NewMockLinear()
			defer mock.Close()
			mock.Handle("IssueUpdate", map[string]any{
				"issueUpdate": map[string]any{"success": true, "issue": map[string]any{"identifier": "ENG-1"}},
			})

			resp, err := linear.IssueClear(context.Background(), mock.Client(), "ENG-1", tc.field)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !resp.IssueUpdate.Success || resp.IssueUpdate.Issue.Identifier != "ENG-1" {
				t.Errorf("response not decoded: %+v", resp.IssueUpdate)
			}

			input, _ := mock.Variables("IssueUpdate")["input"].(map[string]any)
			got, present := input[string(tc.field)]
			if !present {
				t.Fatalf("input %v is missing %q; an omitted key leaves the field unchanged", input, tc.field)
			}
			if !reflect.DeepEqual(got, tc.want) || len(input) != 1 {
				t.Errorf("input = %#v, want only %q: %#v", input, tc.field, tc.want)
			}
		})
	}
}

func TestIssueUpdate_OmitsUnsetFields(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()
	mock.Handle("IssueUpdate", map[string]any{"issueUpdate": map[string]any{"success": true}})

	projectID := "p-uuid"
	if _, err := linear.IssueUpdate(context.Background(), mock.Client(), "ENG-1", linear.IssueUpdateInput{ProjectId: &projectID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := mock.Variables("IssueUpdate")["input"]
	if want := map[string]any{"projectId": "p-uuid"}; !reflect.DeepEqual(input, want) {
		t.Errorf("input = %#v, want %#v", input, want)
	}
}
