package pretty

import "testing"

type stringerVal struct{ s string }

func (v stringerVal) String() string { return v.s }

type namedString string

func TestText(t *testing.T) {
	s := "hi"
	cases := []struct {
		name string
		m    map[string]any
		want string
	}{
		{"absent", map[string]any{}, ""},
		{"nil", map[string]any{"k": nil}, ""},
		{"string", map[string]any{"k": "hi"}, "hi"},
		{"string ptr", map[string]any{"k": &s}, "hi"},
		{"nil string ptr", map[string]any{"k": (*string)(nil)}, ""},
		{"stringer", map[string]any{"k": stringerVal{"who"}}, "who"},
		{"named string falls through", map[string]any{"k": namedString("Active")}, "Active"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Text(tc.m, "k"); got != tc.want {
				t.Errorf("Text = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestStr(t *testing.T) {
	s := "hi"
	if got := Str(map[string]any{"k": "hi"}, "k"); got != "hi" {
		t.Errorf("string: %q", got)
	}
	if got := Str(map[string]any{"k": &s}, "k"); got != "hi" {
		t.Errorf("ptr: %q", got)
	}
	if got := Str(map[string]any{"k": (*string)(nil)}, "k"); got != "" {
		t.Errorf("nil ptr: %q", got)
	}
	if got := Str(map[string]any{}, "k"); got != "" {
		t.Errorf("absent: %q", got)
	}
}

func TestNum(t *testing.T) {
	f := 3.5
	i := 7
	cases := []struct {
		name   string
		v      any
		want   float64
		wantOk bool
	}{
		{"float64", 2.0, 2, true},
		{"float64 ptr", &f, 3.5, true},
		{"int", 4, 4, true},
		{"int ptr", &i, 7, true},
		{"nil float ptr", (*float64)(nil), 0, false},
		{"nil int ptr", (*int)(nil), 0, false},
		{"absent", nil, 0, false},
		{"string", "x", 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := map[string]any{}
			if tc.name != "absent" {
				m["k"] = tc.v
			}
			got, ok := Num(m, "k")
			if got != tc.want || ok != tc.wantOk {
				t.Errorf("Num = (%v, %v), want (%v, %v)", got, ok, tc.want, tc.wantOk)
			}
		})
	}
}

func TestInt(t *testing.T) {
	if got := Int(map[string]any{"k": 3.9}, "k"); got != 3 {
		t.Errorf("truncation: got %d, want 3", got)
	}
	if got := Int(map[string]any{}, "k"); got != 0 {
		t.Errorf("absent: got %d, want 0", got)
	}
}

func TestSubmap(t *testing.T) {
	inner := map[string]any{"a": 1}
	if got := Submap(map[string]any{"k": inner}, "k"); got["a"] != 1 {
		t.Errorf("submap: %v", got)
	}
	if got := Submap(map[string]any{"k": "notamap"}, "k"); got != nil {
		t.Errorf("non-map should be nil, got %v", got)
	}
}

func TestMapSlice(t *testing.T) {
	native := []map[string]any{{"a": 1}}
	if got := MapSlice(map[string]any{"k": native}, "k"); len(got) != 1 || got[0]["a"] != 1 {
		t.Errorf("native shape: %v", got)
	}
	jsonShape := []any{map[string]any{"a": 2}, "skip-non-map"}
	got := MapSlice(map[string]any{"k": jsonShape}, "k")
	if len(got) != 1 || got[0]["a"] != 2 {
		t.Errorf("json shape should keep only maps: %v", got)
	}
	if MapSlice(map[string]any{}, "k") != nil {
		t.Error("absent should be nil")
	}
}

func TestTrimFloat(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{5, "5"},
		{5.5, "5.5"},
		{120, "120"},
		{1000000, "1000000"}, // never scientific notation
	}
	for _, tc := range cases {
		if got := TrimFloat(tc.in); got != tc.want {
			t.Errorf("TrimFloat(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
