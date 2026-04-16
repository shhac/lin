package cycle

import (
	"testing"
	"time"

	"github.com/shhac/lin/internal/linear"
)

func makeCycleNode(id, startsAt, endsAt string) linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle {
	return linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		Id:       id,
		Number:   1,
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}
}

func TestFindNextCycle_Normal(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	nodes := []linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		makeCycleNode("past", "2026-03-01T00:00:00Z", "2026-03-15T00:00:00Z"),
		makeCycleNode("far-future", "2026-06-01T00:00:00Z", "2026-06-15T00:00:00Z"),
		makeCycleNode("near-future", "2026-04-20T00:00:00Z", "2026-05-04T00:00:00Z"),
	}

	got, ok := findNextCycle(nodes, now)
	if !ok {
		t.Fatal("expected to find next cycle")
	}
	if got.Id != "near-future" {
		t.Errorf("Id = %q, want %q", got.Id, "near-future")
	}
}

func TestFindNextCycle_NoFutureCycles(t *testing.T) {
	now := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
	nodes := []linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		makeCycleNode("past-a", "2026-01-01T00:00:00Z", "2026-01-15T00:00:00Z"),
		makeCycleNode("past-b", "2026-06-01T00:00:00Z", "2026-06-15T00:00:00Z"),
	}

	_, ok := findNextCycle(nodes, now)
	if ok {
		t.Error("expected no next cycle")
	}
}

func TestFindNextCycle_SingleCycle(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	nodes := []linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		makeCycleNode("only", "2026-02-01T00:00:00Z", "2026-02-15T00:00:00Z"),
	}

	got, ok := findNextCycle(nodes, now)
	if !ok {
		t.Fatal("expected to find next cycle")
	}
	if got.Id != "only" {
		t.Errorf("Id = %q, want %q", got.Id, "only")
	}
}

func TestFindNextCycle_ParseError(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	nodes := []linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		makeCycleNode("bad", "not-a-date", "also-not-a-date"),
	}

	_, ok := findNextCycle(nodes, now)
	if ok {
		t.Error("expected no result for unparseable dates")
	}
}

func TestFindPreviousCycle_Normal(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	nodes := []linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		makeCycleNode("old", "2026-01-01T00:00:00Z", "2026-01-15T00:00:00Z"),
		makeCycleNode("recent", "2026-03-01T00:00:00Z", "2026-03-15T00:00:00Z"),
		makeCycleNode("future", "2026-05-01T00:00:00Z", "2026-05-15T00:00:00Z"),
	}

	got, ok := findPreviousCycle(nodes, now)
	if !ok {
		t.Fatal("expected to find previous cycle")
	}
	if got.Id != "recent" {
		t.Errorf("Id = %q, want %q", got.Id, "recent")
	}
}

func TestFindPreviousCycle_NoPastCycles(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	nodes := []linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		makeCycleNode("future-a", "2026-01-01T00:00:00Z", "2026-01-15T00:00:00Z"),
		makeCycleNode("future-b", "2026-06-01T00:00:00Z", "2026-06-15T00:00:00Z"),
	}

	_, ok := findPreviousCycle(nodes, now)
	if ok {
		t.Error("expected no previous cycle")
	}
}

func TestFindPreviousCycle_SingleCycle(t *testing.T) {
	now := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	nodes := []linear.TeamCyclesTeamCyclesCycleConnectionNodesCycle{
		makeCycleNode("only", "2026-06-01T00:00:00Z", "2026-06-15T00:00:00Z"),
	}

	got, ok := findPreviousCycle(nodes, now)
	if !ok {
		t.Fatal("expected to find previous cycle")
	}
	if got.Id != "only" {
		t.Errorf("Id = %q, want %q", got.Id, "only")
	}
}
