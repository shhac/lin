package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestReadWrite(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	cfg := &Config{
		DefaultWorkspace: "acme",
		Workspaces: map[string]Workspace{
			"acme": {APIKey: "lin_test_key_abc123", Name: "Acme Corp", URLKey: "acme"},
		},
	}
	if err := Write(cfg); err != nil {
		t.Fatalf("Write: %v", err)
	}

	ClearCache()
	got := Read()
	if got.DefaultWorkspace != "acme" {
		t.Errorf("DefaultWorkspace = %q, want %q", got.DefaultWorkspace, "acme")
	}
	ws, ok := got.Workspaces["acme"]
	if !ok {
		t.Fatal("workspace 'acme' not found")
	}
	if ws.APIKey != "lin_test_key_abc123" {
		t.Errorf("APIKey = %q, want %q", ws.APIKey, "lin_test_key_abc123")
	}
	if ws.Name != "Acme Corp" {
		t.Errorf("Name = %q, want %q", ws.Name, "Acme Corp")
	}
	if ws.URLKey != "acme" {
		t.Errorf("URLKey = %q, want %q", ws.URLKey, "acme")
	}
}

func TestReadDefault(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	got := Read()
	if got.DefaultWorkspace != "" {
		t.Errorf("DefaultWorkspace = %q, want empty", got.DefaultWorkspace)
	}
	if len(got.Workspaces) != 0 {
		t.Errorf("Workspaces length = %d, want 0", len(got.Workspaces))
	}
}

func TestReadCorrupted(t *testing.T) {
	dir := t.TempDir()
	SetConfigDir(dir)
	defer SetConfigDir("")

	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte("{invalid json!!}"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got := Read()
	if got.DefaultWorkspace != "" {
		t.Errorf("DefaultWorkspace = %q, want empty", got.DefaultWorkspace)
	}
	if len(got.Workspaces) != 0 {
		t.Errorf("Workspaces length = %d, want 0", len(got.Workspaces))
	}
}

func TestCacheInvalidation(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	cfg1 := &Config{
		DefaultWorkspace: "first",
		Workspaces:       map[string]Workspace{"first": {APIKey: "key1"}},
	}
	if err := Write(cfg1); err != nil {
		t.Fatalf("Write cfg1: %v", err)
	}

	got1 := Read()
	if got1.DefaultWorkspace != "first" {
		t.Fatalf("expected 'first', got %q", got1.DefaultWorkspace)
	}

	ClearCache()

	cfg2 := &Config{
		DefaultWorkspace: "second",
		Workspaces:       map[string]Workspace{"second": {APIKey: "key2"}},
	}
	if err := Write(cfg2); err != nil {
		t.Fatalf("Write cfg2: %v", err)
	}

	got2 := Read()
	if got2.DefaultWorkspace != "second" {
		t.Errorf("after cache invalidation, DefaultWorkspace = %q, want %q", got2.DefaultWorkspace, "second")
	}
}

func TestConfigDir_XDG(t *testing.T) {
	SetConfigDir("")
	defer SetConfigDir("")

	xdgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdgDir)

	got := ConfigDir()
	want := filepath.Join(xdgDir, "lin")
	if got != want {
		t.Errorf("ConfigDir() = %q, want %q", got, want)
	}
}

// --- workspace.go tests ---

func TestStoreLogin(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	ws := Workspace{APIKey: "lin_test_key_store", Name: "Test Org", URLKey: "testorg"}
	if err := StoreLogin("testorg", ws); err != nil {
		t.Fatalf("StoreLogin: %v", err)
	}

	ClearCache()
	cfg := Read()
	stored, ok := cfg.Workspaces["testorg"]
	if !ok {
		t.Fatal("workspace 'testorg' not found after StoreLogin")
	}
	if stored.Name != "Test Org" {
		t.Errorf("Name = %q, want %q", stored.Name, "Test Org")
	}
}

func TestStoreLogin_SetsDefault(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	ws := Workspace{APIKey: "lin_test_key_default", Name: "First Org"}
	if err := StoreLogin("first", ws); err != nil {
		t.Fatalf("StoreLogin: %v", err)
	}

	cfg := Read()
	if cfg.DefaultWorkspace != "first" {
		t.Errorf("DefaultWorkspace = %q, want %q", cfg.DefaultWorkspace, "first")
	}
}

func TestSetDefaultWorkspace(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	if err := StoreLogin("alpha", Workspace{APIKey: "key_a"}); err != nil {
		t.Fatalf("StoreLogin alpha: %v", err)
	}
	if err := StoreLogin("beta", Workspace{APIKey: "key_b"}); err != nil {
		t.Fatalf("StoreLogin beta: %v", err)
	}

	if err := SetDefaultWorkspace("beta"); err != nil {
		t.Fatalf("SetDefaultWorkspace: %v", err)
	}

	ClearCache()
	if got := GetDefaultWorkspace(); got != "beta" {
		t.Errorf("DefaultWorkspace = %q, want %q", got, "beta")
	}
}

func TestSetDefaultWorkspace_Unknown(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	if err := SetDefaultWorkspace("nonexistent"); err == nil {
		t.Fatal("expected error for unknown workspace")
	}
}

func TestRemoveWorkspace(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	if err := StoreLogin("alpha", Workspace{APIKey: "key_a"}); err != nil {
		t.Fatalf("StoreLogin alpha: %v", err)
	}
	if err := StoreLogin("beta", Workspace{APIKey: "key_b"}); err != nil {
		t.Fatalf("StoreLogin beta: %v", err)
	}
	if err := SetDefaultWorkspace("alpha"); err != nil {
		t.Fatalf("SetDefaultWorkspace: %v", err)
	}

	if err := RemoveWorkspace("alpha"); err != nil {
		t.Fatalf("RemoveWorkspace: %v", err)
	}

	ClearCache()
	cfg := Read()
	if _, ok := cfg.Workspaces["alpha"]; ok {
		t.Error("workspace 'alpha' should be removed")
	}
	if cfg.DefaultWorkspace == "alpha" {
		t.Error("default should rotate away from removed workspace")
	}
	if cfg.DefaultWorkspace != "beta" {
		t.Errorf("default should rotate to 'beta', got %q", cfg.DefaultWorkspace)
	}
}

func TestRemoveWorkspace_Unknown(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	if err := RemoveWorkspace("nonexistent"); err == nil {
		t.Fatal("expected error for unknown workspace")
	}
}

func TestClearAll(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	if err := StoreLogin("ws1", Workspace{APIKey: "key1"}); err != nil {
		t.Fatalf("StoreLogin: %v", err)
	}

	if err := ClearAll(); err != nil {
		t.Fatalf("ClearAll: %v", err)
	}

	ClearCache()
	cfg := Read()
	if len(cfg.Workspaces) != 0 {
		t.Errorf("Workspaces length = %d, want 0", len(cfg.Workspaces))
	}
	if cfg.DefaultWorkspace != "" {
		t.Errorf("DefaultWorkspace = %q, want empty", cfg.DefaultWorkspace)
	}
}

func TestGetWorkspaces(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	if err := StoreLogin("ws1", Workspace{APIKey: "key1", Name: "Org 1"}); err != nil {
		t.Fatalf("StoreLogin ws1: %v", err)
	}
	if err := StoreLogin("ws2", Workspace{APIKey: "key2", Name: "Org 2"}); err != nil {
		t.Fatalf("StoreLogin ws2: %v", err)
	}

	workspaces := GetWorkspaces()
	if len(workspaces) != 2 {
		t.Fatalf("expected 2 workspaces, got %d", len(workspaces))
	}
	if workspaces["ws1"].Name != "Org 1" {
		t.Errorf("ws1 Name = %q, want %q", workspaces["ws1"].Name, "Org 1")
	}
	if workspaces["ws2"].Name != "Org 2" {
		t.Errorf("ws2 Name = %q, want %q", workspaces["ws2"].Name, "Org 2")
	}
}

// --- settings.go tests ---

func TestGetSettings_Empty(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	s := GetSettings()
	if s == nil {
		t.Fatal("GetSettings should never return nil")
	}
	if s.Truncation != nil {
		t.Error("Truncation should be nil for empty settings")
	}
	if s.Pagination != nil {
		t.Error("Pagination should be nil for empty settings")
	}
}

func TestUpdateSettings(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	maxLen := 500
	if err := UpdateSettings(&Settings{
		Truncation: &TruncationSettings{MaxLength: &maxLen},
	}); err != nil {
		t.Fatalf("UpdateSettings: %v", err)
	}

	ClearCache()
	s := GetSettings()
	if s.Truncation == nil || s.Truncation.MaxLength == nil {
		t.Fatal("Truncation.MaxLength should be set")
	}
	if *s.Truncation.MaxLength != 500 {
		t.Errorf("MaxLength = %d, want 500", *s.Truncation.MaxLength)
	}

	// verify settings persist in JSON on disk
	data, err := os.ReadFile(filepath.Join(ConfigDir(), "config.json"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if _, ok := raw["settings"]; !ok {
		t.Error("settings key should be present in config JSON")
	}
}

func TestResetSettings(t *testing.T) {
	SetConfigDir(t.TempDir())
	defer SetConfigDir("")

	maxLen := 100
	if err := UpdateSettings(&Settings{
		Truncation: &TruncationSettings{MaxLength: &maxLen},
	}); err != nil {
		t.Fatalf("UpdateSettings: %v", err)
	}

	if err := ResetSettings(); err != nil {
		t.Fatalf("ResetSettings: %v", err)
	}

	ClearCache()
	s := GetSettings()
	if s.Truncation != nil {
		t.Error("Truncation should be nil after reset")
	}
	if s.Pagination != nil {
		t.Error("Pagination should be nil after reset")
	}
}
