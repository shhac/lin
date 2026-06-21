package issue

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shhac/lin/internal/output/pretty"
)

// renderIssueCard renders a single mapped issue (the map[string]any from
// mappers.MapIssueDetail) as a human-readable card for --format pretty. It reads
// only the fields the IssueGet query already returns; --full sections (relations,
// comment bodies) are layered on separately.
func renderIssueCard(d map[string]any, opts pretty.Options) string {
	c := pretty.New(opts)

	renderTitleLine(c, d, opts)
	c.Rule()
	c.Title(opts.Bold(pretty.Str(d, "title")))
	c.Blank()

	renderMetaGrid(c, d, opts)
	renderDescription(c, d)
	renderAttachments(c, d, opts)
	renderActivity(c, d, opts)
	renderFooter(c, d, opts)

	return c.String()
}

func renderTitleLine(c *pretty.Builder, d map[string]any, opts pretty.Options) {
	id := pretty.Str(d, "identifier")
	status := pretty.Submap(d, "status")
	stateName := pretty.Str(status, "name")
	stateType := pretty.Str(status, "type")

	var parts, plainParts []string
	parts = append(parts, opts.Bold(id), opts.StatusStyle(stateType, stateName))
	plainParts = append(plainParts, id, stateName)

	if prio := pretty.Str(d, "priorityLabel"); prio != "" {
		if p, ok := pretty.Num(d, "priority"); ok && p != 0 {
			parts = append(parts, prio)
			plainParts = append(plainParts, prio)
		}
	}
	if est, ok := pretty.Num(d, "estimate"); ok {
		e := trimFloat(est) + "pts"
		parts = append(parts, e)
		plainParts = append(plainParts, e)
	}

	left := parts[0] + "  " + strings.Join(parts[1:], " · ")
	plainLeft := plainParts[0] + "  " + strings.Join(plainParts[1:], " · ")

	var right, plainRight string
	if rel := opts.RelTime(pretty.Str(d, "updatedAt")); rel != "" {
		plainRight = "updated " + rel
		right = opts.Dim(plainRight)
	}
	c.Header(left, plainLeft, right, plainRight)
}

func renderMetaGrid(c *pretty.Builder, d map[string]any, opts pretty.Options) {
	var pairs [][2]string

	assignee := "Unassigned"
	if a := pretty.Submap(d, "assignee"); a != nil {
		assignee = pretty.Str(a, "name")
	}
	pairs = append(pairs, [2]string{"Assignee", assignee})

	if team := pretty.Submap(d, "team"); team != nil {
		pairs = append(pairs, [2]string{"Team", fmt.Sprintf("%s (%s)", pretty.Str(team, "name"), pretty.Str(team, "key"))})
	}
	if proj := pretty.Submap(d, "project"); proj != nil {
		pairs = append(pairs, [2]string{"Project", pretty.Str(proj, "name")})
	}
	if parent := pretty.Submap(d, "parent"); parent != nil {
		pairs = append(pairs, [2]string{"Parent", pretty.Str(parent, "identifier")})
	}
	if due := pretty.Str(d, "dueDate"); due != "" {
		pairs = append(pairs, [2]string{"Due", pretty.DateOnly(due)})
	}
	if rel := opts.RelTime(pretty.Str(d, "createdAt")); rel != "" {
		pairs = append(pairs, [2]string{"Created", rel})
	}
	c.Grid(pairs)

	if labels := pretty.MapSlice(d, "labels"); len(labels) > 0 {
		names := make([]string, len(labels))
		for i, l := range labels {
			names[i] = pretty.Str(l, "name")
		}
		// Pad "Labels" to the grid's label column so the value aligns with the
		// columns above it.
		labelW := 0
		for _, p := range pairs {
			if n := len(p[0]); n > labelW {
				labelW = n
			}
		}
		c.Field("Labels"+strings.Repeat(" ", labelW-len("Labels")), strings.Join(names, " · "))
	}
}

func renderDescription(c *pretty.Builder, d map[string]any) {
	desc := strings.TrimSpace(pretty.Str(d, "description"))
	if desc == "" {
		return
	}
	c.Blank()
	c.Section("Description")
	c.Wrapped(desc)
}

func renderAttachments(c *pretty.Builder, d map[string]any, opts pretty.Options) {
	atts := pretty.MapSlice(d, "attachments")
	if len(atts) == 0 {
		return
	}
	c.Blank()
	c.Section("Attachments")
	for _, a := range atts {
		src := pretty.Str(a, "sourceType")
		if src == "" {
			src = "link"
		}
		line := opts.Dim(fmt.Sprintf("%-8s", src)) + " " + pretty.Str(a, "title")
		if url := pretty.Str(a, "url"); url != "" {
			line += "  " + opts.Accent(url)
		}
		c.Line(line)
	}
}

func renderActivity(c *pretty.Builder, d map[string]any, opts pretty.Options) {
	comments := pretty.Int(d, "commentCount")
	needs := pretty.Int(d, "customerRequestCount")
	if comments == 0 && needs == 0 {
		return
	}
	c.Blank()
	c.Section("Activity")
	if comments > 0 {
		hint := opts.Dim("→ lin issue comment list " + pretty.Str(d, "identifier"))
		c.Line(fmt.Sprintf("%s   %s", plural(comments, "comment"), hint))
	}
	if needs > 0 {
		important := pretty.Int(d, "customerImportantCount")
		if important > 0 {
			c.Line(fmt.Sprintf("%s (of %d)", plural(important, "important customer need"), needs))
		} else {
			c.Line(plural(needs, "customer need"))
		}
	}
}

func renderFooter(c *pretty.Builder, d map[string]any, opts pretty.Options) {
	c.Blank()
	if branch := pretty.Str(d, "branchName"); branch != "" {
		c.Line(opts.Dim("git branch: ") + branch)
	}
	if url := pretty.Str(d, "url"); url != "" {
		c.Line(opts.Accent(url))
	}
}

func trimFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
