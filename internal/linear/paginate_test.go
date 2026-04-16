package linear

import (
	"fmt"
	"testing"
)

func TestFetchAll_SinglePage(t *testing.T) {
	nodes, err := FetchAll(func(first int, after *string) ([]string, bool, *string, error) {
		return []string{"a", "b", "c"}, false, nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}
	if nodes[0] != "a" || nodes[1] != "b" || nodes[2] != "c" {
		t.Fatalf("unexpected nodes: %v", nodes)
	}
}

func TestFetchAll_MultiplePages(t *testing.T) {
	pages := [][]int{{1, 2}, {3, 4}, {5, 6}}
	callCount := 0
	nodes, err := FetchAll(func(first int, after *string) ([]int, bool, *string, error) {
		page := callCount
		callCount++
		if page == 0 && after != nil {
			t.Fatal("first call should have nil cursor")
		}
		if page > 0 && after == nil {
			t.Fatalf("call %d should have non-nil cursor", page)
		}
		hasNext := page < len(pages)-1
		var cursor *string
		if hasNext {
			c := fmt.Sprintf("cursor-%d", page)
			cursor = &c
		}
		return pages[page], hasNext, cursor, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 6 {
		t.Fatalf("expected 6 nodes, got %d", len(nodes))
	}
	for i, want := range []int{1, 2, 3, 4, 5, 6} {
		if nodes[i] != want {
			t.Fatalf("nodes[%d] = %d, want %d", i, nodes[i], want)
		}
	}
	if callCount != 3 {
		t.Fatalf("expected 3 calls, got %d", callCount)
	}
}

func TestFetchAll_Empty(t *testing.T) {
	nodes, err := FetchAll(func(first int, after *string) ([]string, bool, *string, error) {
		return nil, false, nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nodes) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestFetchAll_ErrorOnFirstPage(t *testing.T) {
	nodes, err := FetchAll(func(first int, after *string) ([]string, bool, *string, error) {
		return nil, false, nil, fmt.Errorf("api error")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "api error" {
		t.Fatalf("unexpected error: %v", err)
	}
	if nodes != nil {
		t.Fatalf("expected nil nodes, got %v", nodes)
	}
}

func TestFetchAll_ErrorOnSecondPage(t *testing.T) {
	callCount := 0
	nodes, err := FetchAll(func(first int, after *string) ([]string, bool, *string, error) {
		callCount++
		if callCount == 1 {
			c := "cursor-1"
			return []string{"a"}, true, &c, nil
		}
		return nil, false, nil, fmt.Errorf("page 2 error")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "page 2 error" {
		t.Fatalf("unexpected error: %v", err)
	}
	if nodes != nil {
		t.Fatalf("expected nil nodes, got %v", nodes)
	}
}
