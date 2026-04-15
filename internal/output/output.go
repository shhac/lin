package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/truncation"
)

const DefaultPageSize = 50

type Pagination struct {
	HasMore    bool   `json:"hasMore"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type PaginatedResult struct {
	Items      any         `json:"items"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func PrintJSON(data any) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	var decoded any
	if err := json.Unmarshal(b, &decoded); err != nil {
		return
	}
	decoded = truncation.Apply(pruneEmpty(decoded))
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	_ = enc.Encode(decoded)
}

func PrintPaginated(items any, pageInfo *Pagination) {
	b, err := json.Marshal(items)
	if err != nil {
		return
	}
	var decoded any
	if err := json.Unmarshal(b, &decoded); err != nil {
		return
	}
	// Prune and truncate each item individually
	if arr, ok := decoded.([]any); ok {
		for i, item := range arr {
			arr[i] = truncation.Apply(pruneEmpty(item))
		}
		decoded = arr
	}
	result := map[string]any{"items": decoded}
	if pageInfo != nil && pageInfo.HasMore {
		result["pagination"] = map[string]any{
			"hasMore":    true,
			"nextCursor": pageInfo.NextCursor,
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	_ = enc.Encode(result)
}

func PrintError(msg string) {
	enc := json.NewEncoder(os.Stderr)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(map[string]string{"error": msg})
	os.Exit(1)
}

func PrintErrorf(format string, args ...any) {
	PrintError(fmt.Sprintf(format, args...))
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

// HandleUnknownCommand registers a handler for unknown subcommands on a cobra command.
func HandleUnknownCommand(cmd *cobra.Command, hint string) {
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			var names []string
			for _, sub := range cmd.Commands() {
				if sub.Name() != "usage" && sub.Name() != "help" {
					names = append(names, sub.Name())
				}
			}
			msg := fmt.Sprintf("Unknown command: %q. Valid commands: %s", args[0], strings.Join(names, ", "))
			if hint != "" {
				msg += ". " + hint
			}
			PrintError(msg)
		}
		return cmd.Help()
	}
}

func pruneEmpty(v any) any {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, v := range val {
			if v == nil {
				continue
			}
			pruned := pruneEmpty(v)
			if isEmpty(pruned) {
				continue
			}
			out[k] = pruned
		}
		if len(out) == 0 {
			return nil
		}
		return out
	case []any:
		out := make([]any, 0, len(val))
		for _, v := range val {
			pruned := pruneEmpty(v)
			if pruned != nil {
				out = append(out, pruned)
			}
		}
		return out
	case string:
		if strings.TrimSpace(val) == "" {
			return nil
		}
		return val
	default:
		return v
	}
}

func isEmpty(v any) bool {
	switch val := v.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(val) == ""
	case map[string]any:
		return len(val) == 0
	case []any:
		return len(val) == 0
	default:
		return false
	}
}
