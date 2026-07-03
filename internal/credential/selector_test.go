package credential

import (
	"runtime"
	"strings"
	"testing"

	"github.com/shhac/lin/internal/config"
)

// stubKeychain swaps the package keychain lookup for the duration of a test so
// the placeholder-hit branch is exercisable without a real keychain.
func stubKeychain(t *testing.T, fn func(string) (string, error)) {
	t.Helper()
	prev := getKeychain
	getKeychain = fn
	t.Cleanup(func() { getKeychain = prev })
}

func placeholderWorkspace() *config.Config {
	return &config.Config{
		DefaultWorkspace: "vault",
		Workspaces: map[string]config.Workspace{
			"vault": {APIKey: keychainPlaceholder, URLKey: "vault"},
		},
	}
}

// seed writes a config directly (plaintext keys, readable on every platform —
// no keychain swap) and points resolution at a scratch dir.
func seed(t *testing.T, cfg *config.Config) {
	t.Helper()
	config.SetConfigDir(t.TempDir())
	t.Cleanup(func() { config.SetConfigDir("") })
	t.Setenv("LINEAR_API_KEY", "")
	t.Cleanup(func() { SetSelectedWorkspace("") })
	if err := config.Write(cfg); err != nil {
		t.Fatalf("write config: %v", err)
	}
	config.ClearCache()
}

func twoWorkspaces() *config.Config {
	return &config.Config{
		DefaultWorkspace: "alpha",
		Workspaces: map[string]config.Workspace{
			"alpha": {APIKey: "lin_alpha_key", URLKey: "alpha"},
			"beta":  {APIKey: "lin_beta_key", URLKey: "beta"},
		},
	}
}

func TestResolveForClient_SelectorResolvesRightWorkspace(t *testing.T) {
	seed(t, twoWorkspaces())
	SetSelectedWorkspace("beta")

	key, err := ResolveForClient()
	if err != nil {
		t.Fatalf("ResolveForClient: %v", err)
	}
	if key != "lin_beta_key" {
		t.Errorf("key = %q, want the beta key (selector must override the default alpha)", key)
	}
}

func TestResolveForClient_UnknownSelector(t *testing.T) {
	seed(t, twoWorkspaces())
	SetSelectedWorkspace("ghost")

	_, err := ResolveForClient()
	if err == nil {
		t.Fatal("unknown selector must error, not fall back")
	}
	if !strings.Contains(err.Error(), "ghost") ||
		!strings.Contains(err.Error(), "alpha") || !strings.Contains(err.Error(), "beta") {
		t.Errorf("error should name the bad alias and list known ones: %v", err)
	}
}

// The gate refuses before ANY fallback — even with a valid default workspace
// AND a LINEAR_API_KEY env value that would otherwise serve the request.
func TestResolveForClient_GateBlocksWithoutSelector(t *testing.T) {
	seed(t, twoWorkspaces())
	t.Setenv("LINEAR_API_KEY", "lin_env_would_serve")
	t.Setenv("LIN_REQUIRE_IDENTITY", "1")

	key, err := ResolveForClient()
	if err == nil {
		t.Fatal("gate must error when set and no --workspace given")
	}
	if key != "" {
		t.Errorf("gate leaked a fallback key %q", key)
	}
	if !strings.Contains(err.Error(), "LIN_REQUIRE_IDENTITY") {
		t.Errorf("error should explain the gate: %v", err)
	}
}

func TestResolveForClient_GateAllowsExplicitSelector(t *testing.T) {
	seed(t, twoWorkspaces())
	t.Setenv("LIN_REQUIRE_IDENTITY", "1")
	SetSelectedWorkspace("beta")

	key, err := ResolveForClient()
	if err != nil {
		t.Fatalf("gate must allow an explicit selector: %v", err)
	}
	if key != "lin_beta_key" {
		t.Errorf("key = %q, want beta", key)
	}
}

func TestResolveForClient_NoGateNoSelectorUsesDefault(t *testing.T) {
	seed(t, twoWorkspaces())

	key, err := ResolveForClient()
	if err != nil {
		t.Fatalf("ResolveForClient: %v", err)
	}
	if key != "lin_alpha_key" {
		t.Errorf("key = %q, want the default alpha key", key)
	}
}

// A placeholder slot whose keychain secret is gone resolves to empty — never a
// fallback to another source. Exercised on non-darwin, where keychain.Get is a
// pure ("", nil); on darwin it would query the real keychain.
func TestResolveForClient_PlaceholderKeychainMiss(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("would query the real macOS keychain; the miss path is covered on other platforms")
	}
	seed(t, placeholderWorkspace())
	SetSelectedWorkspace("vault")

	key, err := ResolveForClient()
	if err != nil {
		t.Fatalf("ResolveForClient: %v", err)
	}
	if key != "" {
		t.Errorf("key = %q, want empty (a placeholder miss must never fall back)", key)
	}
}

// A placeholder slot whose secret is present resolves to the keychain value.
func TestResolveForClient_PlaceholderKeychainHit(t *testing.T) {
	seed(t, placeholderWorkspace())
	stubKeychain(t, func(account string) (string, error) {
		if account != "vault" {
			t.Errorf("keychain queried for %q, want the selected alias vault", account)
		}
		return "lin_vault_secret", nil
	})
	SetSelectedWorkspace("vault")

	key, err := ResolveForClient()
	if err != nil {
		t.Fatalf("ResolveForClient: %v", err)
	}
	if key != "lin_vault_secret" {
		t.Errorf("key = %q, want the keychain secret", key)
	}
}

// workspaceKey backs both resolve() and resolveSelector(); its source strings
// ("keychain"/"config") are what Source() reports, so pin them.
func TestWorkspaceKey_SourceStrings(t *testing.T) {
	stubKeychain(t, func(string) (string, error) { return "sekret", nil })
	if k, src := workspaceKey("vault", config.Workspace{APIKey: keychainPlaceholder}); k != "sekret" || src != "keychain" {
		t.Errorf("placeholder hit = (%q, %q), want (sekret, keychain)", k, src)
	}
	if k, src := workspaceKey("plain", config.Workspace{APIKey: "lin_plain"}); k != "lin_plain" || src != "config" {
		t.Errorf("plaintext = (%q, %q), want (lin_plain, config)", k, src)
	}

	stubKeychain(t, func(string) (string, error) { return "", nil })
	if k, src := workspaceKey("vault", config.Workspace{APIKey: keychainPlaceholder}); k != "" || src != "" {
		t.Errorf("placeholder miss = (%q, %q), want both empty", k, src)
	}
}

func TestRequireIdentity_Values(t *testing.T) {
	cases := map[string]bool{"": false, "0": false, "false": false, "FALSE": false, "1": true, "yes": true}
	for v, want := range cases {
		t.Setenv("LIN_REQUIRE_IDENTITY", v)
		if got := requireIdentity(); got != want {
			t.Errorf("requireIdentity(%q) = %v, want %v", v, got, want)
		}
	}
}
