package resolvers

import (
	"strings"
	"testing"

	"github.com/shhac/lin/internal/testutil"
)

func TestResolveUser_ExactNameMatch(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("UserList", map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Ada Lovelace", "email": "ada@example.com", "displayName": "ada"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Grace Hopper", "email": "grace@example.com", "displayName": "grace"},
				{"id": "cccccccc-1111-2222-3333-444444444444", "name": "Alan Turing", "email": "alan@example.com", "displayName": "alan"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	got, err := ResolveUser(mock.Client(), "Ada Lovelace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "aaaaaaaa-1111-2222-3333-444444444444" {
		t.Errorf("ID = %q", got.ID)
	}
	if got.Name != "Ada Lovelace" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestResolveUser_CaseInsensitive(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("UserList", map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Ada Lovelace", "email": "ada@example.com", "displayName": "ada"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	got, err := ResolveUser(mock.Client(), "ada lovelace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Ada Lovelace" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestResolveUser_EmailMatch(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("UserList", map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Ada Lovelace", "email": "ada@example.com", "displayName": "ada"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Grace Hopper", "email": "grace@example.com", "displayName": "grace"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	got, err := ResolveUser(mock.Client(), "grace@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Email != "grace@example.com" {
		t.Errorf("Email = %q", got.Email)
	}
}

func TestResolveUser_DisplayNameMatch(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("UserList", map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Ada Lovelace", "email": "ada@example.com", "displayName": "babbage-fan"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	got, err := ResolveUser(mock.Client(), "babbage-fan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DisplayName != "babbage-fan" {
		t.Errorf("DisplayName = %q", got.DisplayName)
	}
}

func TestResolveUser_NotFound(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	mock.Handle("UserList", map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Ada Lovelace", "email": "ada@example.com", "displayName": "ada"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveUser(mock.Client(), "Nobody")
	if err == nil {
		t.Fatal("expected error for not found")
	}
	if !strings.Contains(err.Error(), "user not found") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "Ada Lovelace") {
		t.Error("error should list known users")
	}
}

func TestResolveUser_Ambiguous(t *testing.T) {
	mock := testutil.NewMockLinear()
	defer mock.Close()

	// Two users with same display name
	mock.Handle("UserList", map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "aaaaaaaa-1111-2222-3333-444444444444", "name": "Ada Lovelace", "email": "ada1@example.com", "displayName": "ada"},
				{"id": "bbbbbbbb-1111-2222-3333-444444444444", "name": "Ada Smith", "email": "ada2@example.com", "displayName": "ada"},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
		},
	})

	_, err := ResolveUser(mock.Client(), "ada")
	if err == nil {
		t.Fatal("expected error for ambiguous match")
	}
	if !strings.Contains(err.Error(), "ambiguous user") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "2 users") {
		t.Error("error should include match count")
	}
}
