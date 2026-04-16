package mappers

import (
	"testing"

	"github.com/shhac/lin/internal/linear"
)

func TestMapInitiativeSummary_Full(t *testing.T) {
	health := linear.InitiativeUpdateHealthType("onTrack")
	targetDate := "2026-06-30"
	ownerName := "Ada Lovelace"

	input := InitiativeSummaryInput{
		ID:         "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId:     "init-42",
		URL:        "https://linear.app/letsdothis/initiative/init-42",
		Name:       "Platform Reliability",
		Status:     linear.InitiativeStatusActive,
		Health:     &health,
		TargetDate: &targetDate,
		OwnerName:  &ownerName,
	}
	got := MapInitiativeSummary(input)

	assertField(t, got, "id", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	assertField(t, got, "slugId", "init-42")
	assertField(t, got, "url", "https://linear.app/letsdothis/initiative/init-42")
	assertField(t, got, "name", "Platform Reliability")
	if got["status"] != linear.InitiativeStatusActive {
		t.Errorf("status = %v, want %v", got["status"], linear.InitiativeStatusActive)
	}
	if got["health"] != linear.InitiativeUpdateHealthType("onTrack") {
		t.Errorf("health = %v, want onTrack", got["health"])
	}
	assertField(t, got, "targetDate", "2026-06-30")
	if owner, ok := got["owner"].(*string); !ok || *owner != ownerName {
		t.Errorf("owner = %v, want %q", got["owner"], ownerName)
	}
}

func TestMapInitiativeSummary_Minimal(t *testing.T) {
	input := InitiativeSummaryInput{
		ID:     "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId: "init-7",
		URL:    "https://linear.app/letsdothis/initiative/init-7",
		Name:   "Cost Reduction",
		Status: linear.InitiativeStatusPlanned,
	}
	got := MapInitiativeSummary(input)

	if _, ok := got["health"]; ok {
		t.Error("health should be absent when nil")
	}
	if _, ok := got["targetDate"]; ok {
		t.Error("targetDate should be absent when nil")
	}
	if got["owner"] != (*string)(nil) {
		t.Errorf("owner = %v, want nil", got["owner"])
	}
}

func TestFromInitiativeListFields(t *testing.T) {
	health := linear.InitiativeUpdateHealthType("atRisk")
	targetDate := "2026-12-31"

	node := linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative{
		Id:         "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId:     "init-99",
		Url:        "https://linear.app/letsdothis/initiative/init-99",
		Name:       "Data Migration",
		Status:     linear.InitiativeStatusCompleted,
		Health:     &health,
		TargetDate: &targetDate,
		Owner: &linear.InitiativeListInitiativesInitiativeConnectionNodesInitiativeOwnerUser{
			Name: "Grace Hopper",
		},
	}
	got := FromInitiativeListFields(node)

	if got.ID != node.Id {
		t.Errorf("ID = %q, want %q", got.ID, node.Id)
	}
	if got.OwnerName == nil || *got.OwnerName != "Grace Hopper" {
		t.Errorf("OwnerName = %v, want 'Grace Hopper'", got.OwnerName)
	}
}

func TestFromInitiativeListFields_NoOwner(t *testing.T) {
	node := linear.InitiativeListInitiativesInitiativeConnectionNodesInitiative{
		Id:     "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		SlugId: "init-50",
		Url:    "https://linear.app/letsdothis/initiative/init-50",
		Name:   "Unowned Initiative",
		Status: linear.InitiativeStatusPlanned,
	}
	got := FromInitiativeListFields(node)

	if got.OwnerName != nil {
		t.Errorf("OwnerName = %v, want nil", got.OwnerName)
	}
}
