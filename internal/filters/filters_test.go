package filters

import "testing"

func TestBuildIssueFilter_Empty(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{})
	if f != nil {
		t.Fatal("expected nil filter for empty opts")
	}
}

func TestBuildIssueFilter_Team(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Team: "ENG"})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Team == nil {
		t.Fatal("expected team filter")
	}
	if len(f.Team.Or) != 2 {
		t.Fatalf("expected 2 OR branches, got %d", len(f.Team.Or))
	}
	if *f.Team.Or[0].Key.EqIgnoreCase != "ENG" {
		t.Errorf("expected key EqIgnoreCase=ENG, got %s", *f.Team.Or[0].Key.EqIgnoreCase)
	}
	if *f.Team.Or[1].Name.EqIgnoreCase != "ENG" {
		t.Errorf("expected name EqIgnoreCase=ENG, got %s", *f.Team.Or[1].Name.EqIgnoreCase)
	}
}

func TestBuildIssueFilter_AssigneeMe(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Assignee: "me"})
	if f == nil || f.Assignee == nil {
		t.Fatal("expected assignee filter")
	}
	if f.Assignee.IsMe == nil || *f.Assignee.IsMe.Eq != true {
		t.Error("expected IsMe comparator with Eq=true")
	}
}

func TestBuildIssueFilter_AssigneeMe_CaseInsensitive(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Assignee: "ME"})
	if f == nil || f.Assignee == nil || f.Assignee.IsMe == nil {
		t.Fatal("expected IsMe filter for uppercase 'ME'")
	}
}

func TestBuildIssueFilter_AssigneeUUID(t *testing.T) {
	uuid := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	f := BuildIssueFilter(IssueFilterOpts{Assignee: uuid})
	if f == nil || f.Assignee == nil {
		t.Fatal("expected assignee filter")
	}
	if len(f.Assignee.Or) != 4 {
		t.Fatalf("expected 4 OR branches (ID + name + displayName + email), got %d", len(f.Assignee.Or))
	}
	if f.Assignee.Or[0].Id == nil || *f.Assignee.Or[0].Id.Eq != uuid {
		t.Error("expected first branch to be ID match")
	}
}

func TestBuildIssueFilter_AssigneeName(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Assignee: "Ada Lovelace"})
	if f == nil || f.Assignee == nil {
		t.Fatal("expected assignee filter")
	}
	if len(f.Assignee.Or) != 3 {
		t.Fatalf("expected 3 OR branches (name + displayName + email), got %d", len(f.Assignee.Or))
	}
	if *f.Assignee.Or[0].Name.EqIgnoreCase != "Ada Lovelace" {
		t.Error("expected name match")
	}
	if *f.Assignee.Or[1].DisplayName.EqIgnoreCase != "Ada Lovelace" {
		t.Error("expected displayName match")
	}
	if *f.Assignee.Or[2].Email.EqIgnoreCase != "Ada Lovelace" {
		t.Error("expected email match")
	}
}

func TestBuildIssueFilter_Status(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Status: "In Progress"})
	if f == nil || f.State == nil {
		t.Fatal("expected state filter")
	}
	if *f.State.Name.EqIgnoreCase != "In Progress" {
		t.Errorf("expected EqIgnoreCase='In Progress', got %s", *f.State.Name.EqIgnoreCase)
	}
}

func TestBuildIssueFilter_Priority(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Priority: "high"})
	if f == nil || f.Priority == nil {
		t.Fatal("expected priority filter")
	}
	if *f.Priority.Eq != 2.0 {
		t.Errorf("expected priority Eq=2, got %f", *f.Priority.Eq)
	}
}

func TestBuildIssueFilter_InvalidPriority(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Priority: "critical"})
	if f != nil {
		t.Fatal("expected nil filter for invalid priority")
	}
}

func TestBuildIssueFilter_Label(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Label: "bug"})
	if f == nil || f.Labels == nil {
		t.Fatal("expected labels filter")
	}
	if f.Labels.Some == nil || *f.Labels.Some.Name.EqIgnoreCase != "bug" {
		t.Error("expected label name match")
	}
}

func TestBuildIssueFilter_Cycle(t *testing.T) {
	uuid := "c1c2c3c4-d5d6-7890-abcd-ef1234567890"
	f := BuildIssueFilter(IssueFilterOpts{Cycle: uuid})
	if f == nil || f.Cycle == nil {
		t.Fatal("expected cycle filter")
	}
	if *f.Cycle.Id.Eq != uuid {
		t.Errorf("expected cycle ID=%s, got %s", uuid, *f.Cycle.Id.Eq)
	}
}

func TestBuildIssueFilter_DateRanges(t *testing.T) {
	t.Run("updated after and before", func(t *testing.T) {
		f := BuildIssueFilter(IssueFilterOpts{
			UpdatedAfter:  "2025-01-01",
			UpdatedBefore: "2025-06-01",
		})
		if f == nil || f.UpdatedAt == nil {
			t.Fatal("expected updatedAt filter")
		}
		if *f.UpdatedAt.Gte != "2025-01-01" {
			t.Errorf("expected Gte=2025-01-01, got %s", *f.UpdatedAt.Gte)
		}
		if *f.UpdatedAt.Lte != "2025-06-01" {
			t.Errorf("expected Lte=2025-06-01, got %s", *f.UpdatedAt.Lte)
		}
	})

	t.Run("created after only", func(t *testing.T) {
		f := BuildIssueFilter(IssueFilterOpts{CreatedAfter: "2025-03-15"})
		if f == nil || f.CreatedAt == nil {
			t.Fatal("expected createdAt filter")
		}
		if *f.CreatedAt.Gte != "2025-03-15" {
			t.Errorf("expected Gte=2025-03-15, got %s", *f.CreatedAt.Gte)
		}
		if f.CreatedAt.Lte != nil {
			t.Error("expected Lte to be nil")
		}
	})

	t.Run("created before only", func(t *testing.T) {
		f := BuildIssueFilter(IssueFilterOpts{CreatedBefore: "2025-12-31"})
		if f == nil || f.CreatedAt == nil {
			t.Fatal("expected createdAt filter")
		}
		if f.CreatedAt.Gte != nil {
			t.Error("expected Gte to be nil")
		}
		if *f.CreatedAt.Lte != "2025-12-31" {
			t.Errorf("expected Lte=2025-12-31, got %s", *f.CreatedAt.Lte)
		}
	})
}

func TestBuildIssueFilter_MultipleCombined(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{
		Team:     "ENG",
		Status:   "Done",
		Priority: "urgent",
		Assignee: "me",
	})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Team == nil {
		t.Error("expected team filter")
	}
	if f.State == nil {
		t.Error("expected state filter")
	}
	if f.Priority == nil {
		t.Error("expected priority filter")
	}
	if f.Assignee == nil {
		t.Error("expected assignee filter")
	}
}

func TestBuildProjectFilter_UUID(t *testing.T) {
	uuid := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	f := BuildProjectFilter(uuid)
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if len(f.Or) != 3 {
		t.Fatalf("expected 3 OR branches (ID + slugId + name), got %d", len(f.Or))
	}
	if f.Or[0].Id == nil || *f.Or[0].Id.Eq != uuid {
		t.Error("expected first branch to be ID match")
	}
}

func TestBuildProjectFilter_NonUUID(t *testing.T) {
	f := BuildProjectFilter("my-project")
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if len(f.Or) != 2 {
		t.Fatalf("expected 2 OR branches (slugId + name), got %d", len(f.Or))
	}
	if f.Or[0].SlugId == nil || *f.Or[0].SlugId.Eq != "my-project" {
		t.Error("expected slugId match")
	}
	if f.Or[1].Name == nil || *f.Or[1].Name.EqIgnoreCase != "my-project" {
		t.Error("expected name match")
	}
}

func TestBuildNullableProjectFilter_UUID(t *testing.T) {
	uuid := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	f := BuildNullableProjectFilter(uuid)
	if f.Id == nil || *f.Id.Eq != uuid {
		t.Error("expected ID match for UUID input")
	}
}

func TestBuildNullableProjectFilter_NonUUID(t *testing.T) {
	f := BuildNullableProjectFilter("my-project")
	if f.Name == nil || *f.Name.EqIgnoreCase != "my-project" {
		t.Error("expected name match for non-UUID input")
	}
}

func TestBuildTeamFilter(t *testing.T) {
	f := BuildTeamFilter("ENG")
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if len(f.Or) != 2 {
		t.Fatalf("expected 2 OR branches, got %d", len(f.Or))
	}
	if *f.Or[0].Key.EqIgnoreCase != "ENG" {
		t.Error("expected key match")
	}
	if *f.Or[1].Name.EqIgnoreCase != "ENG" {
		t.Error("expected name match")
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"a1b2c3d4-e5f6-7890-abcd-ef1234567890", true},
		{"A1B2C3D4-E5F6-7890-ABCD-EF1234567890", true},
		{"not-a-uuid", false},
		{"", false},
		{"a1b2c3d4-e5f6-7890-abcd", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsUUID(tt.input); got != tt.want {
				t.Errorf("IsUUID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// Verify helper constructors produce the expected comparator shapes.
func TestHelperComparators(t *testing.T) {
	sc := eqIgnoreCase("test")
	if *sc.EqIgnoreCase != "test" {
		t.Error("eqIgnoreCase mismatch")
	}

	ci := containsIgnoreCase("test")
	if *ci.ContainsIgnoreCase != "test" {
		t.Error("containsIgnoreCase mismatch")
	}

	nsc := nullableEqIgnoreCase("test")
	if *nsc.EqIgnoreCase != "test" {
		t.Error("nullableEqIgnoreCase mismatch")
	}
}

// Verify the filter types are correctly shaped (compile-time check + runtime spot check).
func TestBuildIssueFilter_Project_UUID(t *testing.T) {
	uuid := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	f := BuildIssueFilter(IssueFilterOpts{Project: uuid})
	if f == nil || f.Project == nil {
		t.Fatal("expected project filter")
	}
	if f.Project.Id == nil || *f.Project.Id.Eq != uuid {
		t.Error("expected ID match for UUID project input")
	}
}

func TestBuildIssueFilter_Project_Name(t *testing.T) {
	f := BuildIssueFilter(IssueFilterOpts{Project: "my-project"})
	if f == nil || f.Project == nil {
		t.Fatal("expected project filter")
	}
	if f.Project.Name == nil || *f.Project.Name.EqIgnoreCase != "my-project" {
		t.Error("expected name match for non-UUID project input")
	}
}

