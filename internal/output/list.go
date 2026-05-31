package output

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/truncation"
)

const DefaultPageSize = 50

type Pagination struct {
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func PrintList(items any) {
	PrintPaginated(items, nil)
}

func PrintPaginated(items any, pageInfo *Pagination) {
	decoded, ok := cleanList(items)
	if !ok {
		return
	}
	format := ResolveFormat(FormatNDJSON)
	if format == FormatNDJSON {
		printNDJSON(decoded)
		if pageInfo != nil && pageInfo.HasMore {
			_ = NewNDJSONWriter(Stdout()).WritePagination(pageInfo)
		}
		return
	}
	result := map[string]any{"data": decoded}
	if pageInfo != nil && pageInfo.HasMore {
		result["pagination"] = pageInfo
	}
	Print(result, format, false)
}

func cleanList(items any) (any, bool) {
	decoded, ok := toCleanAny(items, false)
	if !ok {
		return nil, false
	}
	if arr, ok := decoded.([]any); ok {
		for i, item := range arr {
			arr[i] = truncation.Apply(pruneEmpty(item))
		}
		return arr, true
	}
	return truncation.Apply(pruneEmpty(decoded)), true
}

// ResolveCursor returns a *string for pagination, or nil if empty.
func ResolveCursor(cursor string) *string {
	if cursor == "" {
		return nil
	}
	return &cursor
}

// ResolvePageSize returns the page size from --limit flag, config, or default.
func ResolvePageSize(limit string) int {
	if limit != "" {
		n, err := strconv.Atoi(limit)
		if err == nil && n > 0 {
			return n
		}
	}
	cfg := config.Read()
	if cfg.Settings != nil && cfg.Settings.Pagination != nil && cfg.Settings.Pagination.DefaultPageSize != nil {
		return *cfg.Settings.Pagination.DefaultPageSize
	}
	return DefaultPageSize
}

// Page holds the bound values of a command's --limit/--cursor flags.
type Page struct {
	limit  string
	cursor string
}

// AddPageFlags registers --limit and --cursor on cmd and returns a Page that
// resolves them lazily via Size() and Cursor().
func AddPageFlags(cmd *cobra.Command) *Page {
	p := &Page{}
	cmd.Flags().StringVar(&p.limit, "limit", "", "Limit results")
	cmd.Flags().StringVar(&p.cursor, "cursor", "", "Pagination cursor for next page")
	return p
}

func (p *Page) Size() int       { return ResolvePageSize(p.limit) }
func (p *Page) Cursor() *string { return ResolveCursor(p.cursor) }

// PrintPage emits a paginated result given items plus genqlient PageInfo fields.
func PrintPage(items any, hasNextPage bool, endCursor *string) {
	cursor := ""
	if endCursor != nil {
		cursor = *endCursor
	}
	PrintPaginated(items, &Pagination{HasMore: hasNextPage, NextCursor: cursor})
}
