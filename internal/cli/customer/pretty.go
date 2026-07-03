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

	right, plainRight := pretty.UpdatedRight(opts, pretty.Str(d, "updatedAt"))
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
		pairs = append(pairs, [2]string{"Revenue", pretty.TrimFloat(r)})
	}
	if s, ok := pretty.Num(d, "size"); ok {
		pairs = append(pairs, [2]string{"Size", pretty.TrimFloat(s)})
	}
	c.Grid(pairs)

	if dom := joinAny(d["domains"]); dom != "" {
		c.Field("Domains", dom)
	}
	if ext := joinAny(d["externalIds"]); ext != "" {
		c.Field("External", ext)
	}

	c.FooterURL(pretty.Str(d, "url"))
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

// joinAny renders a slice of scalars as a " · "-separated list. It accepts the
// mapper's native []string as well as the []any shape a JSON round-trip yields,
// so it works whether the card renders the raw map or a decoded one.
func joinAny(v any) string {
	var parts []string
	switch items := v.(type) {
	case []string:
		parts = items
	case []any:
		parts = make([]string, len(items))
		for i, it := range items {
			parts[i] = fmt.Sprintf("%v", it)
		}
	default:
		return ""
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " · ")
}
