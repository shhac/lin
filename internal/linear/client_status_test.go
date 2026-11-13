package linear

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/shhac/lin/internal/errors"
)

func TestRawQuery_ClassifiesHTTPStatus(t *testing.T) {
	cases := []struct {
		status int
		want   apierrors.FixableBy
	}{
		{http.StatusTooManyRequests, apierrors.FixableByRetry},
		{http.StatusBadGateway, apierrors.FixableByRetry},
		{http.StatusUnauthorized, apierrors.FixableByHuman},
		{http.StatusBadRequest, apierrors.FixableByAgent},
	}
	for _, tc := range cases {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()
			c := NewClient("lin_test_key")
			c.baseURL = srv.URL

			_, err := c.RawQuery(context.Background(), "{ viewer { id } }", nil)

			var apiErr *apierrors.APIError
			if !apierrors.As(err, &apiErr) || apiErr.FixableBy != tc.want {
				t.Errorf("err = %#v, want fixable_by %q", err, tc.want)
			}
		})
	}
}

func TestRawQuery_NetworkFailureIsRetryable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close()
	c := NewClient("lin_test_key")
	c.baseURL = url

	_, err := c.RawQuery(context.Background(), "{ viewer { id } }", nil)

	var apiErr *apierrors.APIError
	if !apierrors.As(err, &apiErr) || apiErr.FixableBy != apierrors.FixableByRetry {
		t.Errorf("err = %#v, want fixable_by retry", err)
	}
}
