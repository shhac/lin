package errors

import (
	stderrors "errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("something went wrong", FixableByHuman)
	if err.Message != "something went wrong" {
		t.Errorf("Message = %q, want %q", err.Message, "something went wrong")
	}
	if err.FixableBy != FixableByHuman {
		t.Errorf("FixableBy = %q, want %q", err.FixableBy, FixableByHuman)
	}
	if err.Error() != "something went wrong" {
		t.Errorf("Error() = %q, want %q", err.Error(), "something went wrong")
	}
}

func TestNewf(t *testing.T) {
	err := Newf(FixableByAgent, "field %q is invalid: got %d", "count", 42)
	want := `field "count" is invalid: got 42`
	if err.Message != want {
		t.Errorf("Message = %q, want %q", err.Message, want)
	}
	if err.FixableBy != FixableByAgent {
		t.Errorf("FixableBy = %q, want %q", err.FixableBy, FixableByAgent)
	}
}

func TestWrap(t *testing.T) {
	cause := fmt.Errorf("underlying failure")
	err := Wrap(cause, FixableByRetry)

	if err.Message != "underlying failure" {
		t.Errorf("Message = %q, want %q", err.Message, "underlying failure")
	}
	if err.FixableBy != FixableByRetry {
		t.Errorf("FixableBy = %q, want %q", err.FixableBy, FixableByRetry)
	}
	if err.Unwrap() != cause {
		t.Error("Unwrap() should return the original cause")
	}
	if !stderrors.Is(err, cause) {
		t.Error("errors.Is should match the wrapped cause")
	}
}

func TestWithHint(t *testing.T) {
	err := New("bad request", FixableByAgent).WithHint("check the input format")
	if err.Hint != "check the input format" {
		t.Errorf("Hint = %q, want %q", err.Hint, "check the input format")
	}
	if err.Message != "bad request" {
		t.Errorf("Message = %q, want %q", err.Message, "bad request")
	}
}

func TestAs(t *testing.T) {
	inner := New("inner error", FixableByHuman).WithHint("fix it")
	wrapped := fmt.Errorf("outer context: %w", inner)

	var target *APIError
	if !As(wrapped, &target) {
		t.Fatal("As should extract APIError from wrapped chain")
	}
	if target.Message != "inner error" {
		t.Errorf("Message = %q, want %q", target.Message, "inner error")
	}
	if target.Hint != "fix it" {
		t.Errorf("Hint = %q, want %q", target.Hint, "fix it")
	}
	if target.FixableBy != FixableByHuman {
		t.Errorf("FixableBy = %q, want %q", target.FixableBy, FixableByHuman)
	}
}

func TestAs_NoMatch(t *testing.T) {
	plain := fmt.Errorf("plain error")
	var target *APIError
	if As(plain, &target) {
		t.Error("As should return false for non-APIError chain")
	}
}
