// Browser credential enrollment for MCP principals: the descriptor rendered by
// lib-agent-mcp's approval flow, and the callback that validates a person's
// pasted Linear API key, stores it under their principal's alias, and returns
// the workspace binding. See lib-agent-mcp design-docs/enrollment.md for the
// trust model; the invariants here are (1) writes are scoped to the verified
// principal's alias, and (2) a slot converges on one Linear organization — a
// key resolving to a different org is an error, never a silent re-point.
package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shhac/lib-agent-mcp/oauth"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/linear"
)

// mcpEnrollmentDescriptor is the single-mode API-key form: Linear has one
// credential shape, so there is no selector to render.
func mcpEnrollmentDescriptor() oauth.CredentialDescriptor {
	return oauth.CredentialDescriptor{
		Title: "Connect your Linear workspace",
		Intro: "One-time setup. Your API key is stored on the server operator's machine and used only for your own calls.",
		Modes: []oauth.CredentialMode{{
			Key: "api-key",
			Fields: []oauth.CredentialField{
				{Key: "api_key", Label: "Linear API key", Secret: true,
					Help: "Create one at Linear → Settings → Security & access → Personal API keys. It looks like lin_api_…"},
			},
		}},
	}
}

// mcpEnroll validates the submitted API key against the Viewer query and stores
// it under alias = principal name. Errors are human-facing form text, not the
// CLI's structured stderr shape. It is an oauth.EnrollFunc, passed directly to
// WithCredentialEnrollment.
func mcpEnroll(ctx context.Context, req oauth.EnrollRequest) (oauth.EnrollResult, error) {
	apiKey := strings.TrimSpace(req.Values["api_key"])
	if apiKey == "" {
		return oauth.EnrollResult{}, errors.New("enter your Linear API key")
	}

	resp, err := linear.Viewer(ctx, linear.ClientWithKey(apiKey))
	if err != nil {
		return oauth.EnrollResult{}, fmt.Errorf("this API key was rejected by Linear: %v", err)
	}
	org := resp.Viewer.Organization
	if org.UrlKey == "" && org.Id == "" {
		return oauth.EnrollResult{}, errors.New(
			"the key was accepted but Linear returned no organization — try a different key")
	}

	if err := checkOrgConvergence(req.Principal, org.Id, org.UrlKey); err != nil {
		return oauth.EnrollResult{}, err
	}
	if err := config.StoreLogin(req.Principal, config.Workspace{
		APIKey: apiKey,
		Name:   org.Name,
		URLKey: org.UrlKey,
		OrgID:  org.Id,
	}); err != nil {
		return oauth.EnrollResult{}, fmt.Errorf("storing the credentials failed: %v", err)
	}
	return oauth.EnrollResult{Binding: map[string]string{bindingKeyWorkspace: req.Principal}}, nil
}

// checkOrgConvergence enforces the one-slot-one-identity rule: strictly by
// alias (never URL/name matching — a principal name must not resolve into
// someone else's workspace record), and only when the stored slot already
// records an organization identity. Prefers the org id (stable across renames)
// and falls back to the urlKey.
func checkOrgConvergence(alias, orgID, urlKey string) error {
	ws, ok := config.GetWorkspaces()[alias]
	if !ok {
		return nil
	}
	switch {
	case ws.OrgID != "" && orgID != "":
		if ws.OrgID != orgID {
			return differentOrgError(alias)
		}
	case ws.URLKey != "":
		if ws.URLKey != urlKey {
			return differentOrgError(alias)
		}
	}
	return nil
}

func differentOrgError(alias string) error {
	return fmt.Errorf(
		"this API key belongs to a different Linear organization than the one enrolled for %q — if that change is intended, ask the operator to re-point your identity (pair add %s --bind workspace=…)",
		alias, alias)
}
