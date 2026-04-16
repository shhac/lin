package mappers

import (
	"testing"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/ptr"
)

func TestMapIssueDetail_Full(t *testing.T) {
	issue := linear.IssueGetIssue{
		Id:            "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Identifier:    "ENG-200",
		Url:           "https://linear.app/test/issue/ENG-200",
		Title:         "Implement caching layer",
		Description:   ptr.To("Add Redis-backed caching for hot paths"),
		BranchName:    "eng-200-caching",
		Priority:      2,
		PriorityLabel: "High",
		Estimate:      ptr.To(5.0),
		DueDate:       ptr.To("2025-06-01"),
		CreatedAt:     "2025-03-01T09:00:00.000Z",
		UpdatedAt:     "2025-04-10T14:30:00.000Z",
		Assignee: &linear.IssueGetIssueAssigneeUser{
			Id:   "11111111-2222-3333-4444-555555555555",
			Name: "Ada Lovelace",
		},
		State: linear.IssueGetIssueStateWorkflowState{
			Id:   "22222222-3333-4444-5555-666666666666",
			Name: "In Progress",
			Type: "started",
		},
		Team: linear.IssueGetIssueTeam{
			Id:   "33333333-4444-5555-6666-777777777777",
			Key:  "ENG",
			Name: "Engineering",
		},
		Project: &linear.IssueGetIssueProject{
			Id:   "44444444-5555-6666-7777-888888888888",
			Name: "Platform Migration",
		},
		Labels: linear.IssueGetIssueLabelsIssueLabelConnection{
			Nodes: []linear.IssueGetIssueLabelsIssueLabelConnectionNodesIssueLabel{
				{Id: "55555555-6666-7777-8888-999999999999", Name: "backend"},
			},
		},
		Parent: &linear.IssueGetIssueParentIssue{
			Id:         "66666666-7777-8888-9999-aaaaaaaaaaaa",
			Identifier: "ENG-100",
		},
	}

	comments := []linear.IssueCommentsIssueCommentsCommentConnectionNodesComment{
		{Id: "c1"}, {Id: "c2"},
	}

	attachments := []linear.IssueAttachmentsIssueAttachmentsAttachmentConnectionNodesAttachment{
		{
			Title:      "PR #42",
			Url:        "https://github.com/org/repo/pull/42",
			SourceType: ptr.To("github"),
		},
	}

	got := MapIssueDetail(issue, comments, attachments)

	assertField(t, got, "id", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assertField(t, got, "identifier", "ENG-200")
	assertField(t, got, "title", "Implement caching layer")
	assertField(t, got, "branchName", "eng-200-caching")

	if got["commentCount"] != 2 {
		t.Errorf("commentCount = %v, want 2", got["commentCount"])
	}

	status := got["status"].(map[string]any)
	if status["name"] != "In Progress" {
		t.Errorf("status.name = %v", status["name"])
	}

	assignee := got["assignee"].(map[string]any)
	if assignee["name"] != "Ada Lovelace" {
		t.Errorf("assignee.name = %v", assignee["name"])
	}

	team := got["team"].(map[string]any)
	if team["key"] != "ENG" {
		t.Errorf("team.key = %v", team["key"])
	}

	project := got["project"].(map[string]any)
	if project["name"] != "Platform Migration" {
		t.Errorf("project.name = %v", project["name"])
	}

	labels := got["labels"].([]map[string]any)
	if len(labels) != 1 || labels[0]["name"] != "backend" {
		t.Errorf("labels = %v", labels)
	}

	atts := got["attachments"].([]map[string]any)
	if len(atts) != 1 || atts[0]["title"] != "PR #42" {
		t.Errorf("attachments = %v", atts)
	}

	parent := got["parent"].(map[string]any)
	if parent["identifier"] != "ENG-100" {
		t.Errorf("parent.identifier = %v", parent["identifier"])
	}
}

func TestMapIssueDetail_NoOptionalFields(t *testing.T) {
	issue := linear.IssueGetIssue{
		Id:            "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Identifier:    "ENG-201",
		Title:         "Minimal issue",
		Priority:      0,
		PriorityLabel: "None",
		State: linear.IssueGetIssueStateWorkflowState{
			Id:   "22222222-3333-4444-5555-666666666666",
			Name: "Backlog",
			Type: "backlog",
		},
		Team: linear.IssueGetIssueTeam{
			Id:   "33333333-4444-5555-6666-777777777777",
			Key:  "ENG",
			Name: "Engineering",
		},
	}

	got := MapIssueDetail(issue, nil, nil)

	if v, ok := got["assignee"].(map[string]any); ok && v != nil {
		t.Errorf("assignee should be nil when unassigned, got %v", v)
	}
	if v, ok := got["project"].(map[string]any); ok && v != nil {
		t.Errorf("project should be nil when unset, got %v", v)
	}
	if v, ok := got["parent"].(map[string]any); ok && v != nil {
		t.Errorf("parent should be nil when unset, got %v", v)
	}
	if got["commentCount"] != 0 {
		t.Errorf("commentCount = %v, want 0", got["commentCount"])
	}
}
