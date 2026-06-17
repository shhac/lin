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

func TestBuildIssueLabelFilter_Empty(t *testing.T) {
	if f := BuildIssueLabelFilter(IssueLabelFilterOpts{}, ""); f != nil {
		t.Errorf("expected nil filter for empty opts, got %+v", f)
	}
}

func TestBuildIssueLabelFilter_Name(t *testing.T) {
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{Name: "Bug"}, "")
	if f == nil || f.Name == nil || f.Name.EqIgnoreCase == nil {
		t.Fatal("expected name EqIgnoreCase filter")
	}
	if *f.Name.EqIgnoreCase != "Bug" {
		t.Errorf("expected EqIgnoreCase=Bug, got %s", *f.Name.EqIgnoreCase)
	}
}

func TestBuildIssueLabelFilter_Search(t *testing.T) {
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{Search: "perf"}, "")
	if f == nil || f.Name == nil || f.Name.ContainsIgnoreCaseAndAccent == nil {
		t.Fatal("expected name ContainsIgnoreCaseAndAccent filter")
	}
	if *f.Name.ContainsIgnoreCaseAndAccent != "perf" {
		t.Errorf("expected ContainsIgnoreCaseAndAccent=perf, got %s", *f.Name.ContainsIgnoreCaseAndAccent)
	}
}

func TestBuildIssueLabelFilter_TeamID(t *testing.T) {
	id := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{}, id)
	if f == nil || f.Team == nil || f.Team.Id == nil {
		t.Fatal("expected team ID filter")
	}
	if *f.Team.Id.Eq != id {
		t.Errorf("expected team ID=%s, got %s", id, *f.Team.Id.Eq)
	}
}

func TestBuildIssueLabelFilter_TeamFlag(t *testing.T) {
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{Team: "ENG"}, "")
	if f == nil || f.Team == nil {
		t.Fatal("expected team filter")
	}
	if len(f.Team.Or) != 2 {
		t.Fatalf("expected 2 OR branches (key + name), got %d", len(f.Team.Or))
	}
	if *f.Team.Or[0].Key.EqIgnoreCase != "ENG" {
		t.Error("expected key match")
	}
}

func TestBuildIssueLabelFilter_IsGroup(t *testing.T) {
	tru := true
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{IsGroup: &tru}, "")
	if f == nil || f.IsGroup == nil || f.IsGroup.Eq == nil {
		t.Fatal("expected isGroup filter")
	}
	if *f.IsGroup.Eq != true {
		t.Errorf("expected isGroup Eq=true, got %v", *f.IsGroup.Eq)
	}
}

func TestBuildIssueLabelFilter_SearchOverwritesName(t *testing.T) {
	// Both opts.Name and opts.Search target f.Name. The current contract: the
	// later assignment (Search) wins. This test pins that behavior.
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{Name: "Bug", Search: "perf"}, "")
	if f == nil || f.Name == nil {
		t.Fatal("expected name filter")
	}
	if f.Name.EqIgnoreCase != nil {
		t.Errorf("expected EqIgnoreCase to be nil when Search is set, got %v", *f.Name.EqIgnoreCase)
	}
	if f.Name.ContainsIgnoreCaseAndAccent == nil || *f.Name.ContainsIgnoreCaseAndAccent != "perf" {
		t.Errorf("expected Search to set ContainsIgnoreCaseAndAccent=perf")
	}
}

func TestBuildIssueLabelFilter_TeamIDOverridesTeamFlag(t *testing.T) {
	id := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{Team: "ENG"}, id)
	if f == nil || f.Team == nil {
		t.Fatal("expected team filter")
	}
	if f.Team.Id == nil || *f.Team.Id.Eq != id {
		t.Errorf("expected resolved teamID to win, got %+v", f.Team)
	}
	if len(f.Team.Or) != 0 {
		t.Errorf("expected no Or branches when teamID is set, got %d", len(f.Team.Or))
	}
}

func TestBuildIssueLabelFilter_AllOptsCombined(t *testing.T) {
	tru := true
	id := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	f := BuildIssueLabelFilter(IssueLabelFilterOpts{Name: "Bug", IsGroup: &tru}, id)
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Name == nil || f.Name.EqIgnoreCase == nil || *f.Name.EqIgnoreCase != "Bug" {
		t.Error("expected name filter")
	}
	if f.Team == nil || f.Team.Id == nil || *f.Team.Id.Eq != id {
		t.Error("expected team filter by ID")
	}
	if f.IsGroup == nil || f.IsGroup.Eq == nil || *f.IsGroup.Eq != true {
		t.Error("expected isGroup filter")
	}
}

func TestBuildProjectLabelFilter_Empty(t *testing.T) {
	if f := BuildProjectLabelFilter(ProjectLabelFilterOpts{}); f != nil {
		t.Errorf("expected nil filter for empty opts, got %+v", f)
	}
}

func TestBuildProjectLabelFilter_Name(t *testing.T) {
	f := BuildProjectLabelFilter(ProjectLabelFilterOpts{Name: "Engineering"})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Name == nil || f.Name.EqIgnoreCase == nil || *f.Name.EqIgnoreCase != "Engineering" {
		t.Errorf("expected EqIgnoreCase=Engineering, got %+v", f.Name)
	}
}

func TestBuildProjectLabelFilter_Search(t *testing.T) {
	f := BuildProjectLabelFilter(ProjectLabelFilterOpts{Search: "eng"})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Name == nil || f.Name.ContainsIgnoreCaseAndAccent == nil || *f.Name.ContainsIgnoreCaseAndAccent != "eng" {
		t.Errorf("expected substring filter, got %+v", f.Name)
	}
}

func TestBuildProjectLabelFilter_IsGroup(t *testing.T) {
	tru := true
	f := BuildProjectLabelFilter(ProjectLabelFilterOpts{IsGroup: &tru})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.IsGroup == nil || f.IsGroup.Eq == nil || *f.IsGroup.Eq != true {
		t.Errorf("expected isGroup=true, got %+v", f.IsGroup)
	}
}

func TestBuildProjectLabelFilter_SearchOverwritesName(t *testing.T) {
	f := BuildProjectLabelFilter(ProjectLabelFilterOpts{Name: "Engineering", Search: "eng"})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Name == nil || f.Name.ContainsIgnoreCaseAndAccent == nil {
		t.Error("expected search to overwrite name")
	}
}

func TestBuildProjectLabelFilter_AllOptsCombined(t *testing.T) {
	tru := true
	f := BuildProjectLabelFilter(ProjectLabelFilterOpts{Name: "Discovery", IsGroup: &tru})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Name == nil || f.Name.EqIgnoreCase == nil || *f.Name.EqIgnoreCase != "Discovery" {
		t.Error("expected name filter")
	}
	if f.IsGroup == nil || f.IsGroup.Eq == nil || *f.IsGroup.Eq != true {
		t.Error("expected isGroup filter")
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

func TestBuildCustomerFilter_Empty(t *testing.T) {
	if f := BuildCustomerFilter(CustomerFilterOpts{}); f != nil {
		t.Fatal("expected nil filter for empty opts")
	}
}

func TestBuildCustomerFilter_TierStatusRevenue(t *testing.T) {
	f := BuildCustomerFilter(CustomerFilterOpts{Tier: "Enterprise", Status: "Active", Revenue: "10000"})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Tier == nil || *f.Tier.DisplayName.EqIgnoreCase != "Enterprise" {
		t.Errorf("tier = %+v", f.Tier)
	}
	if f.Status == nil || *f.Status.Name.EqIgnoreCase != "Active" {
		t.Errorf("status = %+v", f.Status)
	}
	if f.Revenue == nil || f.Revenue.Gte == nil || *f.Revenue.Gte != 10000 {
		t.Errorf("revenue = %+v", f.Revenue)
	}
}

func TestBuildCustomerFilter_SearchAndDomain(t *testing.T) {
	f := BuildCustomerFilter(CustomerFilterOpts{Search: "acme", Domain: "acme.example"})
	if f == nil || f.Name == nil || *f.Name.ContainsIgnoreCaseAndAccent != "acme" {
		t.Errorf("name search = %+v", f.Name)
	}
	if f.Domains == nil || f.Domains.Some == nil || *f.Domains.Some.EqIgnoreCase != "acme.example" {
		t.Errorf("domains = %+v", f.Domains)
	}
}

func TestBuildCustomerNeedFilter_Empty(t *testing.T) {
	if f := BuildCustomerNeedFilter(CustomerNeedFilterOpts{}); f != nil {
		t.Fatal("expected nil filter for empty opts")
	}
}

func TestBuildCustomerNeedFilter_Important(t *testing.T) {
	f := BuildCustomerNeedFilter(CustomerNeedFilterOpts{Important: true})
	if f == nil || f.Priority == nil || f.Priority.Eq == nil || *f.Priority.Eq != 1 {
		t.Errorf("priority = %+v", f)
	}
}

func TestBuildCustomerNeedFilter_UnassignedTriage(t *testing.T) {
	f := BuildCustomerNeedFilter(CustomerNeedFilterOpts{Unassigned: true, Triage: true})
	if f == nil || f.Issue == nil {
		t.Fatal("expected issue sub-filter")
	}
	if f.Issue.Assignee == nil || f.Issue.Assignee.Null == nil || !*f.Issue.Assignee.Null {
		t.Errorf("assignee null = %+v", f.Issue.Assignee)
	}
	if f.Issue.State == nil || f.Issue.State.Type == nil || *f.Issue.State.Type.Eq != "triage" {
		t.Errorf("state type = %+v", f.Issue.State)
	}
}

func TestBuildCustomerNeedFilter_CustomerAndStatus(t *testing.T) {
	f := BuildCustomerNeedFilter(CustomerNeedFilterOpts{Customer: "Acme Corp", Status: "In Progress"})
	if f == nil || f.Customer == nil || *f.Customer.Name.EqIgnoreCase != "Acme Corp" {
		t.Errorf("customer = %+v", f.Customer)
	}
	if f.Issue == nil || f.Issue.State == nil || *f.Issue.State.Name.EqIgnoreCase != "In Progress" {
		t.Errorf("state name = %+v", f.Issue)
	}
}

func TestBuildCustomerFilter_OwnerMe(t *testing.T) {
	f := BuildCustomerFilter(CustomerFilterOpts{Owner: "me"})
	if f == nil || f.Owner == nil {
		t.Fatal("expected owner filter")
	}
	if f.Owner.IsMe == nil || f.Owner.IsMe.Eq == nil || !*f.Owner.IsMe.Eq {
		t.Errorf("expected IsMe comparator with Eq=true, got %+v", f.Owner)
	}
}

func TestBuildCustomerFilter_AllOptsCombined(t *testing.T) {
	f := BuildCustomerFilter(CustomerFilterOpts{
		Tier:    "Enterprise",
		Status:  "Active",
		Owner:   "dana@example.com",
		Domain:  "acme.example",
		Revenue: "25000",
	})
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Tier == nil || f.Status == nil || f.Owner == nil || f.Domains == nil || f.Revenue == nil {
		t.Errorf("expected every facet set, got %+v", f)
	}
	if f.Owner.Or == nil {
		t.Error("expected owner Or branches for a non-me owner")
	}
}

func TestBuildCustomerNeedFilter_ProjectAndDates(t *testing.T) {
	f := BuildCustomerNeedFilter(CustomerNeedFilterOpts{
		Project:       "my-project",
		CreatedAfter:  "2026-01-01",
		CreatedBefore: "2026-02-01",
	})
	if f == nil || f.Project == nil {
		t.Fatal("expected project filter")
	}
	if f.CreatedAt == nil || f.CreatedAt.Gte == nil || *f.CreatedAt.Gte != "2026-01-01" {
		t.Errorf("createdAt gte = %+v", f.CreatedAt)
	}
	if f.CreatedAt.Lte == nil || *f.CreatedAt.Lte != "2026-02-01" {
		t.Errorf("createdAt lte = %+v", f.CreatedAt)
	}
	if f.Issue != nil {
		t.Error("expected no issue sub-filter when only project/dates set")
	}
}
