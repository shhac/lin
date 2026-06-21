package customer

import (
	"fmt"
	"strings"

	"github.com/shhac/lin/internal/output/pretty"
)

// renderCustomerCard renders a customer (the map from MapCustomerDetail) as a
// human-readable card for --format pretty.
func renderCustomerCard(d map[string]any, opts pretty.Options) string {
	c := pretty.New(opts)

	var right, plainRight string
	if rel := opts.RelTime(pretty.Str(d, "updatedAt")); rel != "" {
		plainRight = "updated " + rel
		right = opts.Dim(plainRight)
	}
	name := pretty.Str(d, "name")
	c.Header(opts.Bold(name), name, right, plainRight)
	c.Rule()

	var pairs [][2]string
	if status := pretty.Submap(d, "status"); status != nil {
		name := pretty.Str(status, "name")
		pairs = append(pairs, [2]string{"Status", opts.StatusStyle(customerStatusType(name), name)})
	}
	if owner := pretty.Submap(d, "owner"); owner != nil {
		pairs = append(pairs, [2]string{"Owner", pretty.Str(owner, "name")})
	}
	if tier := pretty.Submap(d, "tier"); tier != nil {
		pairs = append(pairs, [2]string{"Tier", pretty.Str(tier, "displayName")})
	}
	if n := pretty.Int(d, "approximateNeedCount"); n > 0 {
		pairs = append(pairs, [2]string{"Needs", fmt.Sprintf("~%d", n)})
	}
	if r, ok := pretty.Num(d, "revenue"); ok {
		pairs = append(pairs, [2]string{"Revenue", trimFloat(r)})
	}
	if s, ok := pretty.Num(d, "size"); ok {
		pairs = append(pairs, [2]string{"Size", trimFloat(s)})
	}
	c.Grid(pairs)

	if dom := joinAny(d["domains"]); dom != "" {
		c.Field("Domains", dom)
	}
	if ext := joinAny(d["externalIds"]); ext != "" {
		c.Field("External", ext)
	}

	c.Blank()
	if url := pretty.Str(d, "url"); url != "" {
		c.Line(opts.Accent(url))
	}
	return c.String()
}

// customerStatusType maps a workspace-defined customer status name to a
// workflow-state type for coloring. Customer statuses are free-form, so this is
// a best-effort heuristic on common names; unknown names get the neutral color.
func customerStatusType(name string) string {
	switch strings.ToLower(name) {
	case "active", "live", "customer":
		return "completed" // green
	case "churned", "inactive", "lost", "cancelled", "canceled":
		return "canceled" // dim
	case "prospect", "lead", "trial", "onboarding", "evaluating":
		return "started" // yellow
	default:
		return "backlog" // blue
	}
}

// joinAny renders a []any of scalars as a " · "-separated list.
func joinAny(v any) string {
	items, ok := v.([]any)
	if !ok || len(items) == 0 {
		return ""
	}
	parts := make([]string, len(items))
	for i, it := range items {
		parts[i] = fmt.Sprintf("%v", it)
	}
	return strings.Join(parts, " · ")
}

func trimFloat(f float64) string {
	return fmt.Sprintf("%v", f)
}
