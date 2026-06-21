package document

import (
	"strings"

	"github.com/shhac/lin/internal/output/pretty"
)

// renderDocumentCard renders a document (the inline map from the get command) as
// a human-readable card for --format pretty.
func renderDocumentCard(d map[string]any, opts pretty.Options) string {
	c := pretty.New(opts)

	title := pretty.Str(d, "title")
	if icon := pretty.Str(d, "icon"); icon != "" {
		title = icon + " " + title
	}
	var right, plainRight string
	if rel := opts.RelTime(pretty.Str(d, "updatedAt")); rel != "" {
		plainRight = "updated " + rel
		right = opts.Dim(plainRight)
	}
	c.Header(opts.Bold(title), title, right, plainRight)
	c.Rule()

	var pairs [][2]string
	if proj := pretty.Submap(d, "project"); proj != nil {
		pairs = append(pairs, [2]string{"Project", pretty.Str(proj, "name")})
	}
	if creator := pretty.Submap(d, "creator"); creator != nil {
		pairs = append(pairs, [2]string{"Creator", pretty.Str(creator, "name")})
	}
	if up := pretty.Submap(d, "updatedBy"); up != nil {
		pairs = append(pairs, [2]string{"Updated by", pretty.Str(up, "name")})
	}
	if s := pretty.Str(d, "createdAt"); s != "" {
		pairs = append(pairs, [2]string{"Created", pretty.DateOnly(s)})
	}
	c.Grid(pairs)

	if content := strings.TrimSpace(pretty.Str(d, "content")); content != "" {
		c.Blank()
		c.Section("Content")
		c.Wrapped(content)
	}

	c.Blank()
	if url := pretty.Str(d, "url"); url != "" {
		c.Line(opts.Accent(url))
	}
	return c.String()
}
