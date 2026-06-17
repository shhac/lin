package filters

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/priorities"
	"github.com/shhac/lin/internal/ptr"
)

var uuidRE = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// IsUUID returns true if s is a valid UUID v4 format string.
func IsUUID(s string) bool { return uuidRE.MatchString(s) }

// EqIgnoreCase returns a case-insensitive equality StringComparator.
func EqIgnoreCase(s string) *linear.StringComparator {
	return &linear.StringComparator{EqIgnoreCase: ptr.To(s)}
}

// ContainsIgnoreCase returns a case-insensitive substring StringComparator.
func ContainsIgnoreCase(s string) *linear.StringComparator {
	return &linear.StringComparator{ContainsIgnoreCase: ptr.To(s)}
}

// ContainsIgnoreCaseAndAccent returns a case- and accent-insensitive
// substring StringComparator.
func ContainsIgnoreCaseAndAccent(s string) *linear.StringComparator {
	return &linear.StringComparator{ContainsIgnoreCaseAndAccent: ptr.To(s)}
}

// NullableEqIgnoreCase returns a case-insensitive equality NullableStringComparator.
func NullableEqIgnoreCase(s string) *linear.NullableStringComparator {
	return &linear.NullableStringComparator{EqIgnoreCase: ptr.To(s)}
}

// Internal aliases keep package-local call sites concise.
var (
	eqIgnoreCase                = EqIgnoreCase
	containsIgnoreCase          = ContainsIgnoreCase
	containsIgnoreCaseAndAccent = ContainsIgnoreCaseAndAccent
	nullableEqIgnoreCase        = NullableEqIgnoreCase
)

func BuildTeamFilter(input string) *linear.TeamFilter {
	return &linear.TeamFilter{
		Or: []linear.TeamFilter{
			{Key: eqIgnoreCase(input)},
			{Name: eqIgnoreCase(input)},
		},
	}
}

func BuildProjectFilter(input string) *linear.ProjectFilter {
	branches := []linear.ProjectFilter{
		{SlugId: &linear.StringComparator{Eq: ptr.To(input)}},
		{Name: eqIgnoreCase(input)},
	}
	if IsUUID(input) {
		branches = append([]linear.ProjectFilter{{Id: &linear.IDComparator{Eq: ptr.To(input)}}}, branches...)
	}
	return &linear.ProjectFilter{Or: branches}
}

// BuildNullableProjectFilter builds a project filter for use in IssueFilter.Project.
// NullableProjectFilter doesn't support Or, so we pick the best match strategy.
func BuildNullableProjectFilter(input string) *linear.NullableProjectFilter {
	if IsUUID(input) {
		return &linear.NullableProjectFilter{Id: &linear.IDComparator{Eq: ptr.To(input)}}
	}
	return &linear.NullableProjectFilter{
		Name: eqIgnoreCase(input),
	}
}

type IssueLabelFilterOpts struct {
	Name    string // exact match (case-insensitive)
	Search  string // substring match (case+accent insensitive)
	Team    string // team key/name/UUID
	IsGroup *bool
}

func BuildIssueLabelFilter(opts IssueLabelFilterOpts, teamID string) *linear.IssueLabelFilter {
	f := &linear.IssueLabelFilter{}

	if opts.Name != "" {
		f.Name = eqIgnoreCase(opts.Name)
	}
	if opts.Search != "" {
		f.Name = containsIgnoreCaseAndAccent(opts.Search)
	}
	if teamID != "" {
		f.Team = &linear.NullableTeamFilter{Id: &linear.IDComparator{Eq: ptr.To(teamID)}}
	} else if opts.Team != "" {
		f.Team = &linear.NullableTeamFilter{Or: []linear.NullableTeamFilter{
			{Key: eqIgnoreCase(opts.Team)},
			{Name: eqIgnoreCase(opts.Team)},
		}}
	}
	if opts.IsGroup != nil {
		f.IsGroup = &linear.BooleanComparator{Eq: opts.IsGroup}
	}

	if reflect.DeepEqual(*f, linear.IssueLabelFilter{}) {
		return nil
	}
	return f
}

type ProjectLabelFilterOpts struct {
	Name    string // exact match (case-insensitive)
	Search  string // substring match (case+accent insensitive)
	IsGroup *bool
}

func BuildProjectLabelFilter(opts ProjectLabelFilterOpts) *linear.ProjectLabelFilter {
	f := &linear.ProjectLabelFilter{}

	if opts.Name != "" {
		f.Name = eqIgnoreCase(opts.Name)
	}
	if opts.Search != "" {
		f.Name = containsIgnoreCaseAndAccent(opts.Search)
	}
	if opts.IsGroup != nil {
		f.IsGroup = &linear.BooleanComparator{Eq: opts.IsGroup}
	}

	if reflect.DeepEqual(*f, linear.ProjectLabelFilter{}) {
		return nil
	}
	return f
}

type IssueFilterOpts struct {
	Project       string
	Team          string
	Assignee      string
	Status        string
	Priority      string
	Label         string
	Cycle         string
	UpdatedAfter  string
	UpdatedBefore string
	CreatedAfter  string
	CreatedBefore string
}

func BuildIssueFilter(opts IssueFilterOpts) *linear.IssueFilter {
	f := &linear.IssueFilter{}

	if opts.Project != "" {
		f.Project = BuildNullableProjectFilter(opts.Project)
	}

	if opts.Team != "" {
		f.Team = BuildTeamFilter(opts.Team)
	}

	if opts.Assignee != "" {
		f.Assignee = BuildNullableUserFilter(opts.Assignee)
	}

	if opts.Status != "" {
		f.State = &linear.WorkflowStateFilter{Name: eqIgnoreCase(opts.Status)}
	}

	if opts.Priority != "" {
		if p, ok := priorities.Resolve(opts.Priority); ok {
			pf := float64(p)
			f.Priority = &linear.NullableNumberComparator{Eq: &pf}
		}
	}

	if opts.Label != "" {
		f.Labels = &linear.IssueLabelCollectionFilter{
			Some: &linear.IssueLabelFilter{Name: eqIgnoreCase(opts.Label)},
		}
	}

	if opts.Cycle != "" {
		f.Cycle = &linear.NullableCycleFilter{Id: &linear.IDComparator{Eq: ptr.To(opts.Cycle)}}
	}

	if opts.UpdatedAfter != "" || opts.UpdatedBefore != "" {
		dc := &linear.DateComparator{}
		if opts.UpdatedAfter != "" {
			dc.Gte = ptr.To(opts.UpdatedAfter)
		}
		if opts.UpdatedBefore != "" {
			dc.Lte = ptr.To(opts.UpdatedBefore)
		}
		f.UpdatedAt = dc
	}

	if opts.CreatedAfter != "" || opts.CreatedBefore != "" {
		dc := &linear.DateComparator{}
		if opts.CreatedAfter != "" {
			dc.Gte = ptr.To(opts.CreatedAfter)
		}
		if opts.CreatedBefore != "" {
			dc.Lte = ptr.To(opts.CreatedBefore)
		}
		f.CreatedAt = dc
	}

	if reflect.DeepEqual(*f, linear.IssueFilter{}) {
		return nil
	}
	return f
}

// BuildNullableUserFilter builds a user filter accepting "me", a UUID, name,
// display name, or email. Shared by issue assignee and customer owner filters.
func BuildNullableUserFilter(input string) *linear.NullableUserFilter {
	if strings.EqualFold(input, "me") {
		return &linear.NullableUserFilter{IsMe: &linear.BooleanComparator{Eq: ptr.To(true)}}
	}
	branches := []linear.NullableUserFilter{
		{Name: eqIgnoreCase(input)},
		{DisplayName: eqIgnoreCase(input)},
		{Email: eqIgnoreCase(input)},
	}
	if IsUUID(input) {
		branches = append([]linear.NullableUserFilter{{Id: &linear.IDComparator{Eq: ptr.To(input)}}}, branches...)
	}
	return &linear.NullableUserFilter{Or: branches}
}

// BuildNullableCustomerFilter builds a customer filter for use in
// CustomerNeedFilter.Customer, matching a UUID or exact name.
func BuildNullableCustomerFilter(input string) *linear.NullableCustomerFilter {
	if IsUUID(input) {
		return &linear.NullableCustomerFilter{Id: &linear.IDComparator{Eq: ptr.To(input)}}
	}
	return &linear.NullableCustomerFilter{Name: eqIgnoreCase(input)}
}

// BuildCustomerNameFilter builds a customers() filter matching an exact name,
// used as the resolver fallback when an input is not a UUID or slug.
func BuildCustomerNameFilter(input string) *linear.CustomerFilter {
	return &linear.CustomerFilter{Name: eqIgnoreCase(input)}
}

type CustomerFilterOpts struct {
	Search  string // name substring (case+accent insensitive)
	Name    string // name exact (case-insensitive)
	Tier    string // tier display name
	Status  string // status name
	Owner   string // owner: me/name/display name/email/UUID
	Domain  string // email domain (exact)
	Revenue string // minimum revenue (gte)
}

func BuildCustomerFilter(opts CustomerFilterOpts) *linear.CustomerFilter {
	f := &linear.CustomerFilter{}

	if opts.Name != "" {
		f.Name = eqIgnoreCase(opts.Name)
	}
	if opts.Search != "" {
		f.Name = containsIgnoreCaseAndAccent(opts.Search)
	}
	if opts.Tier != "" {
		f.Tier = &linear.CustomerTierFilter{DisplayName: eqIgnoreCase(opts.Tier)}
	}
	if opts.Status != "" {
		f.Status = &linear.CustomerStatusFilter{Name: eqIgnoreCase(opts.Status)}
	}
	if opts.Owner != "" {
		f.Owner = BuildNullableUserFilter(opts.Owner)
	}
	if opts.Domain != "" {
		f.Domains = &linear.StringArrayComparator{Some: &linear.StringItemComparator{EqIgnoreCase: ptr.To(opts.Domain)}}
	}
	if opts.Revenue != "" {
		if n, err := strconv.ParseFloat(opts.Revenue, 64); err == nil {
			f.Revenue = &linear.NumberComparator{Gte: &n}
		}
	}

	if reflect.DeepEqual(*f, linear.CustomerFilter{}) {
		return nil
	}
	return f
}

type CustomerNeedFilterOpts struct {
	Customer      string // customer UUID or name
	Project       string // project UUID, slug, or name
	Important     bool   // priority == 1
	Unassigned    bool   // linked issue has no assignee
	Triage        bool   // linked issue state type == triage
	Status        string // linked issue state name
	Label         string // linked issue label
	Team          string // linked issue team
	CreatedAfter  string
	CreatedBefore string
}

func BuildCustomerNeedFilter(opts CustomerNeedFilterOpts) *linear.CustomerNeedFilter {
	f := &linear.CustomerNeedFilter{}

	if opts.Important {
		f.Priority = &linear.NumberComparator{Eq: ptr.To(1.0)}
	}
	if opts.Customer != "" {
		f.Customer = BuildNullableCustomerFilter(opts.Customer)
	}
	if opts.Project != "" {
		f.Project = BuildNullableProjectFilter(opts.Project)
	}
	if opts.CreatedAfter != "" || opts.CreatedBefore != "" {
		dc := &linear.DateComparator{}
		if opts.CreatedAfter != "" {
			dc.Gte = ptr.To(opts.CreatedAfter)
		}
		if opts.CreatedBefore != "" {
			dc.Lte = ptr.To(opts.CreatedBefore)
		}
		f.CreatedAt = dc
	}

	issue := &linear.NullableIssueFilter{}
	issueSet := false
	if opts.Unassigned {
		issue.Assignee = &linear.NullableUserFilter{Null: ptr.To(true)}
		issueSet = true
	}
	if opts.Triage || opts.Status != "" {
		state := &linear.WorkflowStateFilter{}
		if opts.Triage {
			state.Type = &linear.StringComparator{Eq: ptr.To("triage")}
		}
		if opts.Status != "" {
			state.Name = eqIgnoreCase(opts.Status)
		}
		issue.State = state
		issueSet = true
	}
	if opts.Label != "" {
		issue.Labels = &linear.IssueLabelCollectionFilter{
			Some: &linear.IssueLabelFilter{Name: eqIgnoreCase(opts.Label)},
		}
		issueSet = true
	}
	if opts.Team != "" {
		issue.Team = BuildTeamFilter(opts.Team)
		issueSet = true
	}
	if issueSet {
		f.Issue = issue
	}

	if reflect.DeepEqual(*f, linear.CustomerNeedFilter{}) {
		return nil
	}
	return f
}
