package output

import (
	"testing"

	"github.com/shhac/lin/internal/config"
)

func TestPruneEmpty_NilFields(t *testing.T) {
	input := map[string]any{
		"title": "Test",
		"desc":  nil,
	}
	got := pruneEmpty(input).(map[string]any)
	if _, ok := got["desc"]; ok {
		t.Error("nil field should be pruned")
	}
	if got["title"] != "Test" {
		t.Error("non-nil field should remain")
	}
}

func TestPruneEmpty_EmptyString(t *testing.T) {
	input := map[string]any{
		"name":  "Ada",
		"empty": "",
		"space": "   ",
	}
	got := pruneEmpty(input).(map[string]any)
	if _, ok := got["empty"]; ok {
		t.Error("empty string should be pruned")
	}
	if _, ok := got["space"]; ok {
		t.Error("whitespace-only string should be pruned")
	}
	if got["name"] != "Ada" {
		t.Error("non-empty string should remain")
	}
}

func TestPruneEmpty_EmptyMap(t *testing.T) {
	input := map[string]any{
		"nested": map[string]any{},
		"valid":  map[string]any{"key": "val"},
	}
	got := pruneEmpty(input).(map[string]any)
	if _, ok := got["nested"]; ok {
		t.Error("empty nested map should be pruned")
	}
	nested := got["valid"].(map[string]any)
	if nested["key"] != "val" {
		t.Error("non-empty nested map should remain")
	}
}

func TestPruneEmpty_EmptySlice(t *testing.T) {
	input := map[string]any{
		"items":  []any{},
		"filled": []any{"a", "b"},
	}
	got := pruneEmpty(input).(map[string]any)
	if _, ok := got["items"]; ok {
		t.Error("empty slice should be pruned")
	}
	if filled, ok := got["filled"].([]any); !ok || len(filled) != 2 {
		t.Error("non-empty slice should remain")
	}
}

func TestPruneEmpty_NilValues(t *testing.T) {
	if pruneEmpty(nil) != nil {
		t.Error("nil input should return nil")
	}
}

func TestPruneEmpty_Scalars(t *testing.T) {
	if pruneEmpty(42) != 42 {
		t.Error("integer should pass through")
	}
	if pruneEmpty(true) != true {
		t.Error("bool should pass through")
	}
	if pruneEmpty(3.14) != 3.14 {
		t.Error("float should pass through")
	}
}

func TestPruneEmpty_NestedNilsPruneParent(t *testing.T) {
	input := map[string]any{
		"outer": map[string]any{
			"inner": nil,
		},
	}
	got := pruneEmpty(input)
	if got != nil {
		if m, ok := got.(map[string]any); ok && len(m) > 0 {
			t.Error("map containing only nil nested values should be pruned")
		}
	}
}

func TestPruneEmpty_SliceWithNils(t *testing.T) {
	input := []any{nil, "keep", nil}
	got := pruneEmpty(input).([]any)
	if len(got) != 1 || got[0] != "keep" {
		t.Errorf("expected [keep], got %v", got)
	}
}

func TestResolveCursor_Empty(t *testing.T) {
	if ResolveCursor("") != nil {
		t.Error("empty cursor should return nil")
	}
}

func TestResolveCursor_NonEmpty(t *testing.T) {
	got := ResolveCursor("abc123")
	if got == nil || *got != "abc123" {
		t.Error("non-empty cursor should return pointer to value")
	}
}

func TestResolvePageSize_FromFlag(t *testing.T) {
	// Use a temp config dir so we don't read the user's real config
	tmp := t.TempDir()
	config.SetConfigDir(tmp)
	defer config.SetConfigDir("")

	got := ResolvePageSize("25")
	if got != 25 {
		t.Errorf("expected 25 from flag, got %d", got)
	}
}

func TestResolvePageSize_InvalidFlag(t *testing.T) {
	tmp := t.TempDir()
	config.SetConfigDir(tmp)
	defer config.SetConfigDir("")

	got := ResolvePageSize("not-a-number")
	if got != DefaultPageSize {
		t.Errorf("expected default %d for invalid flag, got %d", DefaultPageSize, got)
	}
}

func TestResolvePageSize_ZeroFlag(t *testing.T) {
	tmp := t.TempDir()
	config.SetConfigDir(tmp)
	defer config.SetConfigDir("")

	got := ResolvePageSize("0")
	if got != DefaultPageSize {
		t.Errorf("expected default %d for zero flag, got %d", DefaultPageSize, got)
	}
}

func TestResolvePageSize_NegativeFlag(t *testing.T) {
	tmp := t.TempDir()
	config.SetConfigDir(tmp)
	defer config.SetConfigDir("")

	got := ResolvePageSize("-5")
	if got != DefaultPageSize {
		t.Errorf("expected default %d for negative flag, got %d", DefaultPageSize, got)
	}
}

func TestResolvePageSize_FromConfig(t *testing.T) {
	tmp := t.TempDir()
	config.SetConfigDir(tmp)
	defer config.SetConfigDir("")

	pageSize := 30
	cfg := &config.Config{
		Workspaces: make(map[string]config.Workspace),
		Settings: &config.Settings{
			Pagination: &config.PaginationSettings{
				DefaultPageSize: &pageSize,
			},
		},
	}
	if err := config.Write(cfg); err != nil {
		t.Fatalf("write config: %v", err)
	}
	config.ClearCache()

	got := ResolvePageSize("")
	if got != 30 {
		t.Errorf("expected 30 from config, got %d", got)
	}
}

func TestResolvePageSize_Fallback(t *testing.T) {
	tmp := t.TempDir()
	config.SetConfigDir(tmp)
	defer config.SetConfigDir("")

	got := ResolvePageSize("")
	if got != DefaultPageSize {
		t.Errorf("expected default %d, got %d", DefaultPageSize, got)
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"whitespace", "  ", true},
		{"non-empty string", "abc", false},
		{"empty map", map[string]any{}, true},
		{"non-empty map", map[string]any{"k": "v"}, false},
		{"empty slice", []any{}, true},
		{"non-empty slice", []any{1}, false},
		{"number", 42, false},
		{"bool", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmpty(tt.v); got != tt.want {
				t.Errorf("isEmpty(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

