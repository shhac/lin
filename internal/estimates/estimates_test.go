package estimates

import (
	"strings"
	"testing"
)

func TestValidEstimates_Exponential(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		cfg := BuildConfig("exponential", false, false)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 4, 8, 16}
		assertIntSlice(t, got, want)
	})

	t.Run("with zero", func(t *testing.T) {
		cfg := BuildConfig("exponential", true, false)
		got := ValidEstimates(cfg)
		want := []int{0, 1, 2, 4, 8, 16}
		assertIntSlice(t, got, want)
	})

	t.Run("extended", func(t *testing.T) {
		cfg := BuildConfig("exponential", false, true)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 4, 8, 16, 32, 64}
		assertIntSlice(t, got, want)
	})

	t.Run("zero and extended", func(t *testing.T) {
		cfg := BuildConfig("exponential", true, true)
		got := ValidEstimates(cfg)
		want := []int{0, 1, 2, 4, 8, 16, 32, 64}
		assertIntSlice(t, got, want)
	})
}

func TestValidEstimates_Fibonacci(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		cfg := BuildConfig("fibonacci", false, false)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 3, 5, 8, 13}
		assertIntSlice(t, got, want)
	})

	t.Run("extended", func(t *testing.T) {
		cfg := BuildConfig("fibonacci", false, true)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 3, 5, 8, 13, 21, 34}
		assertIntSlice(t, got, want)
	})
}

func TestValidEstimates_Linear(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		cfg := BuildConfig("linear", false, false)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 3, 4, 5}
		assertIntSlice(t, got, want)
	})

	t.Run("extended", func(t *testing.T) {
		cfg := BuildConfig("linear", false, true)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		assertIntSlice(t, got, want)
	})
}

func TestValidEstimates_TShirt(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		cfg := BuildConfig("tShirt", false, false)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 3, 4, 5}
		assertIntSlice(t, got, want)
	})

	t.Run("extended", func(t *testing.T) {
		cfg := BuildConfig("tShirt", false, true)
		got := ValidEstimates(cfg)
		want := []int{1, 2, 3, 4, 5, 6}
		assertIntSlice(t, got, want)
	})
}

func TestValidEstimates_UnknownScale(t *testing.T) {
	cfg := BuildConfig("unknown", false, false)
	got := ValidEstimates(cfg)
	if got != nil {
		t.Errorf("expected nil for unknown scale, got %v", got)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := BuildConfig("fibonacci", false, false)
	for _, v := range []int{1, 2, 3, 5, 8, 13} {
		if err := Validate(cfg, v); err != nil {
			t.Errorf("Validate(%d) unexpected error: %v", v, err)
		}
	}
}

func TestValidate_Invalid(t *testing.T) {
	cfg := BuildConfig("fibonacci", false, false)
	err := Validate(cfg, 7)
	if err == nil {
		t.Fatal("expected error for invalid estimate")
	}
	if !strings.Contains(err.Error(), "invalid estimate 7") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_NotUsed(t *testing.T) {
	cfg := BuildConfig("notUsed", false, false)
	err := Validate(cfg, 1)
	if err == nil {
		t.Fatal("expected error for notUsed team")
	}
	if !strings.Contains(err.Error(), "does not use estimates") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_ZeroNotAllowed(t *testing.T) {
	cfg := BuildConfig("linear", false, false)
	err := Validate(cfg, 0)
	if err == nil {
		t.Fatal("expected error for zero when not allowed")
	}
}

func TestValidate_ZeroAllowed(t *testing.T) {
	cfg := BuildConfig("linear", true, false)
	if err := Validate(cfg, 0); err != nil {
		t.Errorf("unexpected error when zero is allowed: %v", err)
	}
}

func TestFormatScale_Numeric(t *testing.T) {
	got := FormatScale("fibonacci", []int{1, 2, 3, 5, 8})
	want := "1 | 2 | 3 | 5 | 8"
	if got != want {
		t.Errorf("FormatScale = %q, want %q", got, want)
	}
}

func TestFormatScale_TShirt(t *testing.T) {
	got := FormatScale("tShirt", []int{1, 2, 3, 4, 5, 6})
	if !strings.Contains(got, "(XS)") {
		t.Errorf("expected t-shirt label XS in %q", got)
	}
	if !strings.Contains(got, "(XXL)") {
		t.Errorf("expected t-shirt label XXL in %q", got)
	}
	want := "1 (XS) | 2 (S) | 3 (M) | 4 (L) | 5 (XL) | 6 (XXL)"
	if got != want {
		t.Errorf("FormatScale = %q, want %q", got, want)
	}
}

func TestFormatScale_TShirt_WithZero(t *testing.T) {
	got := FormatScale("tShirt", []int{0, 1, 2})
	if !strings.Contains(got, "0 (None)") {
		t.Errorf("expected t-shirt label None for 0 in %q", got)
	}
}

func assertIntSlice(t *testing.T, got, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %v (len %d), want %v (len %d)", got, len(got), want, len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d: got %d, want %d", i, got[i], want[i])
		}
	}
}
