package filters

import (
	"regexp"
	"strings"

	"github.com/shhac/lin/internal/linear"
	"github.com/shhac/lin/internal/priorities"
)

var uuidRE = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func isUUID(s string) bool { return uuidRE.MatchString(s) }

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
func intPtr(i int) *int       { return &i }

func eqIgnoreCase(s string) *linear.StringComparator {
	return &linear.StringComparator{EqIgnoreCase: strPtr(s)}
}

func containsIgnoreCase(s string) *linear.StringComparator {
	return &linear.StringComparator{ContainsIgnoreCase: strPtr(s)}
}

func nullableEqIgnoreCase(s string) *linear.NullableStringComparator {
	return &linear.NullableStringComparator{EqIgnoreCase: strPtr(s)}
}

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
		{SlugId: &linear.StringComparator{Eq: strPtr(input)}},
		{Name: eqIgnoreCase(input)},
	}
	if isUUID(input) {
		branches = append([]linear.ProjectFilter{{Id: &linear.IDComparator{Eq: strPtr(input)}}}, branches...)
	}
	return &linear.ProjectFilter{Or: branches}
}

// BuildNullableProjectFilter builds a project filter for use in IssueFilter.Project.
// NullableProjectFilter doesn't support Or, so we pick the best match strategy.
func BuildNullableProjectFilter(input string) *linear.NullableProjectFilter {
	if isUUID(input) {
		return &linear.NullableProjectFilter{Id: &linear.IDComparator{Eq: strPtr(input)}}
	}
	return &linear.NullableProjectFilter{
		Name: eqIgnoreCase(input),
	}
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
	empty := true

	if opts.Project != "" {
		f.Project = BuildNullableProjectFilter(opts.Project)
		empty = false
	}

	if opts.Team != "" {
		f.Team = BuildTeamFilter(opts.Team)
		empty = false
	}

	if opts.Assignee != "" {
		lower := strings.ToLower(opts.Assignee)
		if lower == "me" {
			f.Assignee = &linear.NullableUserFilter{IsMe: &linear.BooleanComparator{Eq: boolPtr(true)}}
		} else {
			branches := []linear.NullableUserFilter{
				{Name: eqIgnoreCase(opts.Assignee)},
				{DisplayName: eqIgnoreCase(opts.Assignee)},
				{Email: eqIgnoreCase(opts.Assignee)},
			}
			if isUUID(opts.Assignee) {
				branches = append([]linear.NullableUserFilter{{Id: &linear.IDComparator{Eq: strPtr(opts.Assignee)}}}, branches...)
			}
			f.Assignee = &linear.NullableUserFilter{Or: branches}
		}
		empty = false
	}

	if opts.Status != "" {
		f.State = &linear.WorkflowStateFilter{Name: eqIgnoreCase(opts.Status)}
		empty = false
	}

	if opts.Priority != "" {
		if p, ok := priorities.Resolve(opts.Priority); ok {
			pf := float64(p)
			f.Priority = &linear.NullableNumberComparator{Eq: &pf}
			empty = false
		}
	}

	if opts.Label != "" {
		f.Labels = &linear.IssueLabelCollectionFilter{
			Some: &linear.IssueLabelFilter{Name: eqIgnoreCase(opts.Label)},
		}
		empty = false
	}

	if opts.Cycle != "" {
		f.Cycle = &linear.NullableCycleFilter{Id: &linear.IDComparator{Eq: strPtr(opts.Cycle)}}
		empty = false
	}

	if opts.UpdatedAfter != "" || opts.UpdatedBefore != "" {
		dc := &linear.DateComparator{}
		if opts.UpdatedAfter != "" {
			dc.Gte = strPtr(opts.UpdatedAfter)
		}
		if opts.UpdatedBefore != "" {
			dc.Lte = strPtr(opts.UpdatedBefore)
		}
		f.UpdatedAt = dc
		empty = false
	}

	if opts.CreatedAfter != "" || opts.CreatedBefore != "" {
		dc := &linear.DateComparator{}
		if opts.CreatedAfter != "" {
			dc.Gte = strPtr(opts.CreatedAfter)
		}
		if opts.CreatedBefore != "" {
			dc.Lte = strPtr(opts.CreatedBefore)
		}
		f.CreatedAt = dc
		empty = false
	}

	if empty {
		return nil
	}
	return f
}
