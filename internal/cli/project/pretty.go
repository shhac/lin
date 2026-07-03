package project

import (
	"fmt"
	"strings"

	"github.com/shhac/lin/internal/output/pretty"
)

// renderProjectCard renders a project (the inline map from the get command) as a
// human-readable card for --format pretty.
func renderProjectCard(d map[string]any, opts pretty.Options) string {
	c := pretty.New(opts)

	var right, plainRight string
	if p, ok := pretty.Num(d, "progress"); ok {
		plainRight = fmt.Sprintf("%.0f%% complete", p*100)
		right = opts.Dim(plainRight)
	}
	name := pretty.Str(d, "name")
	c.Header(opts.Bold(name), name, right, plainRight)
	c.Rule()

	if desc := strings.TrimSpace(pretty.Str(d, "description")); desc != "" {
		c.Wrapped(desc)
		c.Blank()
	}

	var pairs [][2]string
	if state := pretty.Text(d, "status"); state != "" {
		pairs = append(pairs, [2]string{"Status", opts.StatusStyle(state, pretty.Capitalize(state))})
	}
	if lead := pretty.Submap(d, "lead"); lead != nil {
		pairs = append(pairs, [2]string{"Lead", pretty.Str(lead, "name")})
	}
	if s := pretty.Str(d, "startDate"); s != "" {
		pairs = append(pairs, [2]string{"Start", pretty.DateOnly(s)})
	}
	if s := pretty.Str(d, "targetDate"); s != "" {
		pairs = append(pairs, [2]string{"Target", pretty.DateOnly(s)})
	}
	c.Grid(pairs)
	if labels := labelNames(d); labels != "" {
		c.Field("Labels", labels)
	}

	if content := strings.TrimSpace(pretty.Str(d, "content")); content != "" {
		c.Blank()
		c.Section("Content")
		c.Wrapped(content)
	}

	if ms := pretty.MapSlice(d, "milestones"); len(ms) > 0 {
		c.Blank()
		c.Section("Milestones")
		for _, m := range ms {
			line := "• " + pretty.Str(m, "name")
			if t := pretty.Str(m, "targetDate"); t != "" {
				line += opts.Dim("  (" + pretty.DateOnly(t) + ")")
			}
			c.Line(line)
		}
	}

	c.FooterURL(pretty.Str(d, "url"))
	return c.String()
}

func labelNames(d map[string]any) string {
	labels := pretty.MapSlice(d, "labels")
	if len(labels) == 0 {
		return ""
	}
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = pretty.Str(l, "name")
	}
	return strings.Join(names, " · ")
}
