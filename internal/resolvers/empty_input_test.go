package resolvers

import (
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"

	apierrors "github.com/shhac/lin/internal/errors"
	"github.com/shhac/lin/internal/testutil"
)

// The mock registers no handlers, so any request a resolver made would fail
// with "no handler" instead of the empty-input error: passing proves blank
// input is rejected before it can reach Linear's match-all comparators.
func TestResolvers_RejectBlankInput(t *testing.T) {
	resolvers := []struct {
		noun    string
		resolve func(graphql.Client, string) error
	}{
		{"project", func(c graphql.Client, in string) error { _, err := ResolveProject(c, in); return err }},
		{"team", func(c graphql.Client, in string) error { _, err := ResolveTeam(c, in); return err }},
		{"customer", func(c graphql.Client, in string) error { _, err := ResolveCustomer(c, in); return err }},
		{"document", func(c graphql.Client, in string) error { _, err := ResolveDocument(c, in); return err }},
		{"initiative", func(c graphql.Client, in string) error { _, err := ResolveInitiative(c, in); return err }},
		{"user", func(c graphql.Client, in string) error { _, err := ResolveUser(c, in); return err }},
		{"status", func(c graphql.Client, in string) error { _, err := ResolveWorkflowState(c, in, "team-uuid"); return err }},
		{"labels", func(c graphql.Client, in string) error { _, err := ResolveIssueLabels(c, in, ""); return err }},
		{"labels", func(c graphql.Client, in string) error { _, err := ResolveProjectLabels(c, in); return err }},
	}
	inputs := []string{"", "   ", "\t\n"}

	mock := testutil.NewMockLinear()
	defer mock.Close()

	for _, r := range resolvers {
		for _, in := range inputs {
			t.Run(r.noun+"/"+strings.ReplaceAll(in, " ", "·"), func(t *testing.T) {
				err := r.resolve(mock.Client(), in)
				assertEmptyInputErr(t, err, r.noun)
			})
		}
	}
}

func TestLabelResolvers_RejectOnlySeparators(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	_, err := ResolveIssueLabels(mock.Client(), " , ,", "")
	assertEmptyInputErr(t, err, "labels")

	_, err = ResolveProjectLabels(mock.Client(), ",,")
	assertEmptyInputErr(t, err, "labels")
}

func assertEmptyInputErr(t *testing.T, err error, noun string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error for blank input")
	}
	if want := noun + " must not be empty"; err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByAgent {
		t.Errorf("expected FixableByAgent APIError, got %#v", err)
	}
}
