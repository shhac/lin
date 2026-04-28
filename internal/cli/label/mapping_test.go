package label

import (
	"testing"

	"github.com/shhac/lin/internal/linear"
)

func TestMapLabel_MinimalFields(t *testing.T) {
	got := mapIssueLabel(linear.IssueLabelFields{
		Id:    "label-1",
		Name:  "Bug",
		Color: "#ff0000",
	})
	if got["id"] != "label-1" {
		t.Errorf("id = %v", got["id"])
	}
	if got["name"] != "Bug" {
		t.Errorf("name = %v", got["name"])
	}
	if got["color"] != "#ff0000" {
		t.Errorf("color = %v", got["color"])
	}
	if _, ok := got["description"]; ok {
		t.Error("description should be absent for nil pointer")
	}
	if _, ok := got["isGroup"]; ok {
		t.Error("isGroup should be absent when false")
	}
	if _, ok := got["team"]; ok {
		t.Error("team should be absent for nil pointer")
	}
	if _, ok := got["parent"]; ok {
		t.Error("parent should be absent for nil pointer")
	}
}

func TestMapLabel_EmptyDescriptionString(t *testing.T) {
	empty := ""
	got := mapIssueLabel(linear.IssueLabelFields{Id: "x", Name: "y", Color: "#000", Description: &empty})
	if _, ok := got["description"]; ok {
		t.Error("empty-string description should be omitted")
	}
}

func TestMapLabel_PopulatedDescription(t *testing.T) {
	desc := "Tests added or improved"
	got := mapIssueLabel(linear.IssueLabelFields{Id: "x", Name: "y", Color: "#000", Description: &desc})
	if got["description"] != desc {
		t.Errorf("description = %v", got["description"])
	}
}

func TestMapLabel_GroupLabel(t *testing.T) {
	got := mapIssueLabel(linear.IssueLabelFields{Id: "x", Name: "y", Color: "#000", IsGroup: true})
	if got["isGroup"] != true {
		t.Errorf("isGroup = %v", got["isGroup"])
	}
}

func TestMapProjectLabel_MinimalFields(t *testing.T) {
	got := mapProjectLabel(linear.ProjectLabelFields{
		Id:    "plabel-1",
		Name:  "Discovery",
		Color: "#abcdef",
	})
	if got["id"] != "plabel-1" {
		t.Errorf("id = %v", got["id"])
	}
	if got["name"] != "Discovery" {
		t.Errorf("name = %v", got["name"])
	}
	if got["color"] != "#abcdef" {
		t.Errorf("color = %v", got["color"])
	}
	if _, ok := got["team"]; ok {
		t.Error("project labels should never expose a team field")
	}
}

func TestMapProjectLabel_GroupAndParent(t *testing.T) {
	got := mapProjectLabel(linear.ProjectLabelFields{
		Id:      "plabel-1",
		Name:    "Quality",
		Color:   "#000",
		IsGroup: true,
		Parent: &linear.ProjectLabelFieldsParentProjectLabel{
			Id:   "parent-1",
			Name: "Quality Group",
		},
	})
	if got["isGroup"] != true {
		t.Errorf("isGroup = %v", got["isGroup"])
	}
	parent, ok := got["parent"].(map[string]any)
	if !ok {
		t.Fatalf("parent not a map: %T", got["parent"])
	}
	if parent["name"] != "Quality Group" {
		t.Errorf("parent.name = %v", parent["name"])
	}
}

func TestMapLabel_TeamAndParent(t *testing.T) {
	got := mapIssueLabel(linear.IssueLabelFields{
		Id:    "x",
		Name:  "y",
		Color: "#000",
		Team: &linear.IssueLabelFieldsTeam{
			Id:   "team-1",
			Key:  "ENG",
			Name: "Engineering",
		},
		Parent: &linear.IssueLabelFieldsParentIssueLabel{
			Id:   "parent-1",
			Name: "Quality",
		},
	})
	team, ok := got["team"].(map[string]any)
	if !ok {
		t.Fatalf("team not a map: %T", got["team"])
	}
	if team["key"] != "ENG" {
		t.Errorf("team.key = %v", team["key"])
	}
	if team["name"] != "Engineering" {
		t.Errorf("team.name = %v", team["name"])
	}
	parent, ok := got["parent"].(map[string]any)
	if !ok {
		t.Fatalf("parent not a map: %T", got["parent"])
	}
	if parent["name"] != "Quality" {
		t.Errorf("parent.name = %v", parent["name"])
	}
}
