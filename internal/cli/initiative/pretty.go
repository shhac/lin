package initiative

import (
	"strings"

	"github.com/shhac/lin/internal/output/pretty"
)

// renderInitiativeCard renders an initiative (the inline map from the get
// command) as a human-readable card for --format pretty.
func renderInitiativeCard(d map[string]any, opts pretty.Options) string {
	c := pretty.New(opts)

	right, plainRight := pretty.UpdatedRight(opts, pretty.Str(d, "updatedAt"))
	name := pretty.Str(d, "name")
	c.Header(opts.Bold(name), name, right, plainRight)
	c.Rule()

	var pairs [][2]string
	if s := pretty.Text(d, "status"); s != "" {
		pairs = append(pairs, [2]string{"Status", opts.StatusStyle(initiativeStatusType(s), s)})
	}
	if owner := pretty.Submap(d, "owner"); owner != nil {
		pairs = append(pairs, [2]string{"Owner", pretty.Str(owner, "name")})
	}
	if h := pretty.Text(d, "health"); h != "" {
		pairs = append(pairs, [2]string{"Health", opts.HealthStyle(h)})
	}
	if s := pretty.Str(d, "targetDate"); s != "" {
		pairs = append(pairs, [2]string{"Target", pretty.DateOnly(s)})
	}
	if s := pretty.Str(d, "startedAt"); s != "" {
		pairs = append(pairs, [2]string{"Started", pretty.DateOnly(s)})
	}
	if s := pretty.Str(d, "completedAt"); s != "" {
		pairs = append(pairs, [2]string{"Completed", pretty.DateOnly(s)})
	}
	c.Grid(pairs)

	if desc := strings.TrimSpace(pretty.Str(d, "description")); desc != "" {
		c.Blank()
		c.Section("Description")
		c.Wrapped(desc)
	}
	if content := strings.TrimSpace(pretty.Str(d, "content")); content != "" {
		c.Blank()
		c.Section("Content")
		c.Wrapped(content)
	}

	c.FooterURL(pretty.Str(d, "url"))
	return c.String()
}

// initiativeStatusType maps an initiative lifecycle status (Planned/Active/
// Completed) to the workflow-state type StatusStyle colors by.
func initiativeStatusType(status string) string {
	switch status {
	case "Active":
		return "started"
	case "Completed":
		return "completed"
	default: // Planned and any unknown
		return "backlog"
	}
}
