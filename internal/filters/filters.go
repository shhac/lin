package filters

import (
	"reflect"
	"regexp"
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
	Name   string // exact match (case-insensitive)
	Search string // substring match (case+accent insensitive)
	Team   string // team key/name/UUID
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

type IssueFilterOpts struct {
	Project      string
	Team         string
	Assignee     string
	Status       string
	Priority     string
	Label        string
	Cycle        string
	UpdatedAfter string
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
		lower := strings.ToLower(opts.Assignee)
		if lower == "me" {
			f.Assignee = &linear.NullableUserFilter{IsMe: &linear.BooleanComparator{Eq: ptr.To(true)}}
		} else {
			branches := []linear.NullableUserFilter{
				{Name: eqIgnoreCase(opts.Assignee)},
				{DisplayName: eqIgnoreCase(opts.Assignee)},
				{Email: eqIgnoreCase(opts.Assignee)},
			}
			if IsUUID(opts.Assignee) {
				branches = append([]linear.NullableUserFilter{{Id: &linear.IDComparator{Eq: ptr.To(opts.Assignee)}}}, branches...)
			}
			f.Assignee = &linear.NullableUserFilter{Or: branches}
		}
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
