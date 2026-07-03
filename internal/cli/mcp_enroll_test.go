package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/shhac/lib-agent-mcp/oauth"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/linear"
)

// fakeLinear answers the Viewer GraphQL POST with queued response bodies (FIFO
// with a sticky tail, so the convergence test can hand out a foreign identity
// once and the matching one thereafter) and records the Authorization header
// each call carried.
type fakeLinear struct {
	mu        sync.Mutex
	responses []string
	tokens    []string
}

func (f *fakeLinear) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tokens = append(f.tokens, r.Header.Get("Authorization"))
	body := `{"errors":[{"message":"no response queued"}]}`
	if len(f.responses) > 0 {
		body = f.responses[0]
		if len(f.responses) > 1 {
			f.responses = f.responses[1:]
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, body)
}

func (f *fakeLinear) queue(bodies ...string) { f.responses = append(f.responses, bodies...) }

func viewerBody(orgID, urlKey, name string) string {
	return fmt.Sprintf(
		`{"data":{"viewer":{"id":"u1","name":"Alice","email":"a@x.com","displayName":"alice","organization":{"id":%q,"name":%q,"urlKey":%q}}}}`,
		orgID, name, urlKey)
}

type enrollFixture struct {
	fake   *fakeLinear
	enroll oauth.EnrollFunc
}

func newEnrollFixture(t *testing.T) *enrollFixture {
	t.Helper()
	// Force StoreLogin's file fallback (darwin) and never touch a real keychain;
	// isolate config to a scratch dir; keep env from serving the request.
	t.Setenv("LIN_NO_KEYCHAIN", "1")
	t.Setenv("LINEAR_API_KEY", "")
	config.SetConfigDir(t.TempDir())
	t.Cleanup(func() { config.SetConfigDir("") })

	fake := &fakeLinear{}
	ts := httptest.NewServer(fake)
	t.Cleanup(ts.Close)
	linear.Configure(linear.Options{BaseURL: ts.URL})
	t.Cleanup(func() { linear.Configure(linear.Options{}) })

	return &enrollFixture{fake: fake, enroll: mcpEnroll}
}

func (f *enrollFixture) run(t *testing.T, principal string, values map[string]string) (oauth.EnrollResult, error) {
	t.Helper()
	return f.enroll(context.Background(), oauth.EnrollRequest{Principal: principal, Mode: "api-key", Values: values})
}

func TestMCPEnrollHappyPath(t *testing.T) {
	f := newEnrollFixture(t)
	f.fake.queue(viewerBody("org-acme", "acme", "Acme Inc"))

	res, err := f.run(t, "alice", map[string]string{"api_key": "lin_api_alice"})
	if err != nil {
		t.Fatalf("enroll: %v", err)
	}
	if res.Binding["workspace"] != "alice" {
		t.Errorf("binding = %v, want workspace=alice", res.Binding)
	}
	// The submitted key was the one validated against Linear.
	if len(f.fake.tokens) != 1 || f.fake.tokens[0] != "lin_api_alice" {
		t.Errorf("auth headers = %v, want one call with the submitted key", f.fake.tokens)
	}

	ws, ok := config.GetWorkspaces()["alice"]
	if !ok {
		t.Fatal("no workspace stored under the principal alias")
	}
	// Alias = principal name; org metadata derived from the Viewer response.
	if ws.Name != "Acme Inc" || ws.URLKey != "acme" || ws.OrgID != "org-acme" {
		t.Errorf("stored workspace metadata = %+v", ws)
	}
	if config.GetDefaultWorkspace() != "alice" {
		t.Errorf("first enrollee should become default, got %q", config.GetDefaultWorkspace())
	}
	// The exact stored key is only observable where the keychain opt-out forces
	// the plaintext file fallback (darwin); elsewhere it is the placeholder.
	if runtime.GOOS == "darwin" && ws.APIKey != "lin_api_alice" {
		t.Errorf("stored APIKey = %q, want the submitted key", ws.APIKey)
	}
}

func TestMCPEnrollRejectedKey(t *testing.T) {
	f := newEnrollFixture(t)
	f.fake.queue(`{"errors":[{"message":"authentication failed - invalid api key"}]}`)

	_, err := f.run(t, "alice", map[string]string{"api_key": "lin_api_bad"})
	if err == nil || !strings.Contains(err.Error(), "rejected") {
		t.Errorf("err = %v, want the rejection surfaced", err)
	}
	if _, ok := config.GetWorkspaces()["alice"]; ok {
		t.Error("rejected key must not be stored")
	}
}

func TestMCPEnrollEmptyKey(t *testing.T) {
	f := newEnrollFixture(t)
	if _, err := f.run(t, "alice", map[string]string{"api_key": "  "}); err == nil {
		t.Error("blank key must error before calling Linear")
	}
	if len(f.fake.tokens) != 0 {
		t.Errorf("Linear called %d times before local validation", len(f.fake.tokens))
	}
}

// The convergence rule: a slot that already records an org identity only accepts
// a key proving the same org.
func TestMCPEnrollConvergence(t *testing.T) {
	f := newEnrollFixture(t)
	// Pre-store alice bound to org-acme (written directly so the pre-state is
	// deterministic regardless of keychain).
	if err := config.Write(&config.Config{
		DefaultWorkspace: "alice",
		Workspaces: map[string]config.Workspace{
			"alice": {APIKey: "lin_api_old", Name: "Acme Inc", URLKey: "acme", OrgID: "org-acme"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	config.ClearCache()

	// First call resolves a different org, later calls the matching one.
	f.fake.queue(
		viewerBody("org-other", "other", "Other Corp"),
		viewerBody("org-acme", "acme", "Acme Inc"),
	)

	// Different org → refused, slot untouched.
	if _, err := f.run(t, "alice", map[string]string{"api_key": "lin_api_bob"}); err == nil ||
		!strings.Contains(err.Error(), "different Linear organization") {
		t.Errorf("err = %v, want a different-org refusal", err)
	}
	if ws := config.GetWorkspaces()["alice"]; ws.OrgID != "org-acme" || ws.Name != "Acme Inc" {
		t.Errorf("slot mutated by refused enrollment: %+v", ws)
	}

	config.ClearCache()
	// Same org → idempotent re-enroll succeeds and returns the binding.
	res, err := f.run(t, "alice", map[string]string{"api_key": "lin_api_new"})
	if err != nil {
		t.Fatalf("re-enroll same org: %v", err)
	}
	if res.Binding["workspace"] != "alice" {
		t.Errorf("binding = %v", res.Binding)
	}
	if runtime.GOOS == "darwin" {
		if ws := config.GetWorkspaces()["alice"]; ws.APIKey != "lin_api_new" {
			t.Errorf("re-enroll did not update the key in place: %+v", ws)
		}
	}
}

// A principal whose name matches another workspace's urlKey/name must not
// converge against it — the check is strictly by alias.
func TestMCPEnrollConvergenceIsAliasStrict(t *testing.T) {
	f := newEnrollFixture(t)
	if err := config.Write(&config.Config{
		DefaultWorkspace: "ops",
		Workspaces: map[string]config.Workspace{
			"ops": {APIKey: "lin_api_ops", Name: "acme", URLKey: "acme", OrgID: "org-ops"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	config.ClearCache()
	f.fake.queue(viewerBody("org-acme", "acme", "Acme Inc"))

	if _, err := f.run(t, "acme", map[string]string{"api_key": "lin_api_acme"}); err != nil {
		t.Fatalf("enrollment blocked by an unrelated workspace matching the principal: %v", err)
	}
}

// The legacy-user gate: a slot that recorded only a urlKey (no OrgID, as older
// logins did) still converges — by urlKey — so a foreign key can't re-point it.
func TestMCPEnrollConvergenceUrlKeyFallback(t *testing.T) {
	f := newEnrollFixture(t)
	// Legacy slot: URLKey known, OrgID never recorded — the fallback anchor.
	if err := config.Write(&config.Config{
		DefaultWorkspace: "alice",
		Workspaces: map[string]config.Workspace{
			"alice": {APIKey: "lin_api_old", Name: "Acme Inc", URLKey: "acme"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	config.ClearCache()

	// First call resolves a different urlKey, later calls the matching one.
	f.fake.queue(
		viewerBody("org-other", "other", "Other Corp"),
		viewerBody("org-acme", "acme", "Acme Inc"),
	)

	// Different urlKey → refused, slot untouched (still legacy, no OrgID).
	if _, err := f.run(t, "alice", map[string]string{"api_key": "lin_api_bob"}); err == nil ||
		!strings.Contains(err.Error(), "different Linear organization") {
		t.Errorf("err = %v, want a different-org refusal via the urlKey fallback", err)
	}
	if ws := config.GetWorkspaces()["alice"]; ws.URLKey != "acme" || ws.OrgID != "" || ws.APIKey != "lin_api_old" {
		t.Errorf("slot mutated by refused enrollment: %+v", ws)
	}

	config.ClearCache()
	// Same urlKey → succeeds and now records the OrgID for stronger future checks.
	res, err := f.run(t, "alice", map[string]string{"api_key": "lin_api_new"})
	if err != nil {
		t.Fatalf("same-urlKey re-enroll: %v", err)
	}
	if res.Binding["workspace"] != "alice" {
		t.Errorf("binding = %v", res.Binding)
	}
	ws := config.GetWorkspaces()["alice"]
	if ws.OrgID != "org-acme" {
		t.Errorf("successful re-enroll should record the OrgID going forward: %+v", ws)
	}
	if runtime.GOOS == "darwin" && ws.APIKey != "lin_api_new" {
		t.Errorf("re-enroll did not update the key in place: %+v", ws)
	}
}

// A key Linear accepts but that resolves to no organization is refused, and
// nothing is stored.
func TestMCPEnrollNoOrganization(t *testing.T) {
	f := newEnrollFixture(t)
	f.fake.queue(viewerBody("", "", ""))

	_, err := f.run(t, "alice", map[string]string{"api_key": "lin_api_alice"})
	if err == nil || !strings.Contains(err.Error(), "no organization") {
		t.Errorf("err = %v, want the no-organization refusal", err)
	}
	if _, ok := config.GetWorkspaces()["alice"]; ok {
		t.Error("a key resolving to no organization must not be stored")
	}
}

func TestMCPEnrollmentDescriptorShape(t *testing.T) {
	d := mcpEnrollmentDescriptor()
	if len(d.Modes) != 1 || d.Modes[0].Key != "api-key" {
		t.Fatalf("modes = %+v, want the single api-key mode", d.Modes)
	}
	fields := d.Modes[0].Fields
	if len(fields) != 1 || fields[0].Key != "api_key" {
		t.Fatalf("fields = %+v, want a single api_key field", fields)
	}
	if fields[0].Optional || !fields[0].Secret {
		t.Errorf("api_key must be a required secret: %+v", fields[0])
	}
}

func TestMCPIdentityBinding(t *testing.T) {
	// Bound principal → --workspace selector plus the fail-closed gate.
	argv, env := mcpIdentityBinding(oauth.Verified{
		PrincipalGrant: oauth.PrincipalGrant{Name: "alice", Binding: map[string]string{"workspace": "alice"}},
	})
	if strings.Join(argv, " ") != "--workspace alice" {
		t.Errorf("argv = %v, want --workspace alice", argv)
	}
	if strings.Join(env, ",") != "LIN_REQUIRE_IDENTITY=1" {
		t.Errorf("env = %v, want the identity gate", env)
	}

	// Unbound principal → gate only, no selector (fail closed).
	argv, env = mcpIdentityBinding(oauth.Verified{
		PrincipalGrant: oauth.PrincipalGrant{Name: "bob"},
	})
	if argv != nil {
		t.Errorf("argv = %v, want nil for an unbound principal", argv)
	}
	if strings.Join(env, ",") != "LIN_REQUIRE_IDENTITY=1" {
		t.Errorf("env = %v, want the identity gate", env)
	}
}
