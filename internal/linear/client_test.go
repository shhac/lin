package linear

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRawQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "lin_test_key" {
			t.Errorf("expected Authorization header 'lin_test_key', got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got %q", r.Header.Get("Content-Type"))
		}

		var req graphQLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Query != "{ viewer { id } }" {
			t.Errorf("unexpected query: %s", req.Query)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(graphQLResponse{
			Data: json.RawMessage(`{"viewer":{"id":"user-123"}}`),
		})
	}))
	defer srv.Close()

	c := NewClient("lin_test_key")
	c.http = srv.Client()
	// Point at test server instead of real API
	origURL := defaultAPIURL
	defer func() { setAPIURL(origURL) }()
	setAPIURL(srv.URL)

	data, err := c.RawQuery(context.Background(), "{ viewer { id } }", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result struct {
		Viewer struct {
			ID string `json:"id"`
		} `json:"viewer"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.Viewer.ID != "user-123" {
		t.Errorf("expected user-123, got %s", result.Viewer.ID)
	}
}

func TestGraphQLErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(graphQLResponse{
			Errors: []graphQLError{{Message: "Entity not found"}},
		})
	}))
	defer srv.Close()

	c := NewClient("lin_test_key")
	c.http = srv.Client()
	origURL := defaultAPIURL
	defer func() { setAPIURL(origURL) }()
	setAPIURL(srv.URL)

	_, err := c.RawQuery(context.Background(), "{ issue(id: \"bad\") { id } }", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "graphql: Entity not found" {
		t.Errorf("unexpected error: %v", err)
	}
}
