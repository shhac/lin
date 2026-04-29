package credential

import (
	"testing"

	"github.com/shhac/lin/internal/config"
)

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
