package credential

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shhac/lin/internal/config"
)

// TestStoreLogin_Headless_FileFallback exercises the real credential-WRITE path
// (config.StoreLogin) non-interactively. Setting the per-CLI keychain opt-out
// (derived by lib-agent-cli from the "app.paulie.lin" service) makes
// creds.Keychain.Available() report false, so keychain.Store returns
// ErrKeychainUnavailable and StoreLogin deterministically keeps the plaintext
// key in the config file instead of swapping in the __KEYCHAIN__ placeholder —
// on every platform, including darwin, where it would otherwise reach the
// `security` CLI and its GUI prompt. Using LIN_NO_KEYCHAIN (not the family-wide
// LIB_AGENT_NO_KEYCHAIN) also proves the lib's prefix derivation.
func TestStoreLogin_Headless_FileFallback(t *testing.T) {
	// On non-darwin platforms keychain.Store is a no-op that succeeds without
	// touching the backend, so StoreLogin always writes the placeholder and the
	// opt-out has nothing to bypass. The interesting path — opt-out forcing the
	// file fallback that would otherwise hit the macOS `security` GUI prompt —
	// only exists on darwin.
	if runtime.GOOS != "darwin" {
		t.Skip("keychain write path is darwin-only")
	}
	t.Setenv("LIN_NO_KEYCHAIN", "1")
	t.Setenv("LINEAR_API_KEY", "")
	dir := t.TempDir()
	config.SetConfigDir(dir)
	t.Cleanup(func() { config.SetConfigDir("") })

	const key = "lin_headless_write_key"
	if err := config.StoreLogin("headlessws", config.Workspace{APIKey: key, Name: "Headless"}); err != nil {
		t.Fatalf("StoreLogin: %v", err)
	}

	// Under the keychain opt-out the secret must land in the config file as
	// plaintext, NOT as the __KEYCHAIN__ placeholder.
	raw, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatalf("config file not written: %v", err)
	}
	if !strings.Contains(string(raw), key) {
		t.Errorf("plaintext key not written to config under keychain opt-out:\n%s", raw)
	}
	if strings.Contains(string(raw), "__KEYCHAIN__") {
		t.Errorf("keychain placeholder written despite opt-out:\n%s", raw)
	}
	var onDisk config.Config
	if err := json.Unmarshal(raw, &onDisk); err != nil {
		t.Fatalf("config not valid JSON: %v", err)
	}
	if got := onDisk.Workspaces["headlessws"].APIKey; got != key {
		t.Errorf("stored APIKey = %q, want plaintext %q (file fallback)", got, key)
	}

	// Round-trip via the read path: resolved from config, not keychain.
	config.ClearCache()
	if got := Resolve(); got != key {
		t.Errorf("Resolve() = %q, want %q", got, key)
	}
	if src := Source(); src != "config" {
		t.Errorf("Source() = %q, want \"config\" (keychain opt-out forces file fallback)", src)
	}

	// Remove and confirm it's gone.
	if err := config.RemoveWorkspace("headlessws"); err != nil {
		t.Fatalf("RemoveWorkspace: %v", err)
	}
	config.ClearCache()
	if _, ok := config.GetWorkspaces()["headlessws"]; ok {
		t.Error("workspace still present after RemoveWorkspace")
	}
	if got := Resolve(); got != "" {
		t.Errorf("Resolve() after remove = %q, want empty", got)
	}
}

func TestResolve_EnvVar(t *testing.T) {
	config.SetConfigDir(t.TempDir())
	defer config.SetConfigDir("")

	t.Setenv("LINEAR_API_KEY", "lin_env_test_key_abc")

	got := Resolve()
	if got != "lin_env_test_key_abc" {
		t.Errorf("Resolve() = %q, want %q", got, "lin_env_test_key_abc")
	}

	src := Source()
	if src != "environment" {
		t.Errorf("Source() = %q, want %q", src, "environment")
	}
}

func TestResolve_ConfigPlaintext(t *testing.T) {
	config.SetConfigDir(t.TempDir())
	defer config.SetConfigDir("")

	t.Setenv("LINEAR_API_KEY", "")

	cfg := &config.Config{
		DefaultWorkspace: "testws",
		Workspaces: map[string]config.Workspace{
			"testws": {APIKey: "lin_config_test_key_xyz", Name: "Test WS"},
		},
	}
	if err := config.Write(cfg); err != nil {
		t.Fatalf("Write: %v", err)
	}
	config.ClearCache()

	got := Resolve()
	if got != "lin_config_test_key_xyz" {
		t.Errorf("Resolve() = %q, want %q", got, "lin_config_test_key_xyz")
	}

	src := Source()
	if src != "config" {
		t.Errorf("Source() = %q, want %q", src, "config")
	}
}

func TestResolve_LegacyApiKey(t *testing.T) {
	config.SetConfigDir(t.TempDir())
	defer config.SetConfigDir("")

	t.Setenv("LINEAR_API_KEY", "")

	cfg := &config.Config{
		LegacyAPIKey: "lin_legacy_key_abc",
		Workspaces:   map[string]config.Workspace{},
	}
	if err := config.Write(cfg); err != nil {
		t.Fatalf("Write: %v", err)
	}
	config.ClearCache()

	got := Resolve()
	if got != "lin_legacy_key_abc" {
		t.Errorf("Resolve() = %q, want %q", got, "lin_legacy_key_abc")
	}

	src := Source()
	if src != "config" {
		t.Errorf("Source() = %q, want %q", src, "config")
	}
}

func TestResolve_NoKey(t *testing.T) {
	config.SetConfigDir(t.TempDir())
	defer config.SetConfigDir("")

	t.Setenv("LINEAR_API_KEY", "")

	got := Resolve()
	if got != "" {
		t.Errorf("Resolve() = %q, want empty", got)
	}

	src := Source()
	if src != "" {
		t.Errorf("Source() = %q, want empty", src)
	}
}

func TestResolve_EnvOverridesConfig(t *testing.T) {
	config.SetConfigDir(t.TempDir())
	defer config.SetConfigDir("")

	t.Setenv("LINEAR_API_KEY", "lin_env_wins")

	cfg := &config.Config{
		DefaultWorkspace: "testws",
		Workspaces: map[string]config.Workspace{
			"testws": {APIKey: "lin_config_loses", Name: "Test WS"},
		},
	}
	if err := config.Write(cfg); err != nil {
		t.Fatalf("Write: %v", err)
	}
	config.ClearCache()

	got := Resolve()
	if got != "lin_env_wins" {
		t.Errorf("Resolve() = %q, want %q", got, "lin_env_wins")
	}

	src := Source()
	if src != "environment" {
		t.Errorf("Source() = %q, want %q", src, "environment")
	}
}

func TestSource(t *testing.T) {
	tests := []struct {
		name    string
		envKey  string
		cfg     *config.Config
		wantSrc string
	}{
		{
			name:    "environment source",
			envKey:  "lin_env_key",
			cfg:     &config.Config{Workspaces: map[string]config.Workspace{}},
			wantSrc: "environment",
		},
		{
			name:   "config source via workspace",
			envKey: "",
			cfg: &config.Config{
				DefaultWorkspace: "ws",
				Workspaces:       map[string]config.Workspace{"ws": {APIKey: "lin_key"}},
			},
			wantSrc: "config",
		},
		{
			name:    "config source via legacy key",
			envKey:  "",
			cfg:     &config.Config{LegacyAPIKey: "lin_legacy", Workspaces: map[string]config.Workspace{}},
			wantSrc: "config",
		},
		{
			name:    "no source",
			envKey:  "",
			cfg:     &config.Config{Workspaces: map[string]config.Workspace{}},
			wantSrc: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SetConfigDir(t.TempDir())
			defer config.SetConfigDir("")
			t.Setenv("LINEAR_API_KEY", tt.envKey)

			if err := config.Write(tt.cfg); err != nil {
				t.Fatalf("Write: %v", err)
			}
			config.ClearCache()

			got := Source()
			if got != tt.wantSrc {
				t.Errorf("Source() = %q, want %q", got, tt.wantSrc)
			}
		})
	}
}
