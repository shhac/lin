package truncation

import (
	"strings"
	"testing"
)

func resetState() {
	Configure(ConfigOpts{})
}

func TestApply_ShortContent(t *testing.T) {
	resetState()
	input := map[string]any{
		"title":       "Test Issue",
		"description": "Short text",
	}
	got := Apply(input).(map[string]any)
	if got["description"] != "Short text" {
		t.Errorf("expected short content unchanged, got %v", got["description"])
	}
	if got["descriptionLength"] != 10 {
		t.Errorf("expected descriptionLength=10, got %v", got["descriptionLength"])
	}
}

func TestApply_LongDescription(t *testing.T) {
	resetState()
	long := strings.Repeat("x", 300)
	input := map[string]any{
		"description": long,
	}
	got := Apply(input).(map[string]any)
	desc := got["description"].(string)
	if len(desc) > defaultMaxLength+len(ellipsis)+5 {
		t.Errorf("expected truncated description, got length %d", len(desc))
	}
	if !strings.HasSuffix(desc, ellipsis) {
		t.Error("expected ellipsis suffix")
	}
	if got["descriptionLength"] != 300 {
		t.Errorf("expected descriptionLength=300, got %v", got["descriptionLength"])
	}
}

func TestApply_CustomMaxLength(t *testing.T) {
	Configure(ConfigOpts{MaxLength: 50})
	defer resetState()

	long := strings.Repeat("y", 100)
	input := map[string]any{
		"description": long,
	}
	got := Apply(input).(map[string]any)
	desc := got["description"].(string)
	// 50 chars + ellipsis
	if len(desc) > 55 {
		t.Errorf("expected truncated at ~50 chars, got length %d", len(desc))
	}
}

func TestApply_FullMode(t *testing.T) {
	Configure(ConfigOpts{Full: true})
	defer resetState()

	long := strings.Repeat("z", 500)
	input := map[string]any{
		"description": long,
	}
	got := Apply(input).(map[string]any)
	if got["description"] != long {
		t.Error("expected full content in --full mode")
	}
}

func TestApply_ExpandSpecificField(t *testing.T) {
	Configure(ConfigOpts{Expand: "description"})
	defer resetState()

	longDesc := strings.Repeat("a", 500)
	longBody := strings.Repeat("b", 500)
	input := map[string]any{
		"description": longDesc,
		"body":        longBody,
	}
	got := Apply(input).(map[string]any)
	if got["description"] != longDesc {
		t.Error("expected description expanded")
	}
	body := got["body"].(string)
	if !strings.HasSuffix(body, ellipsis) {
		t.Error("expected body still truncated")
	}
}

func TestApply_NonTruncatableField(t *testing.T) {
	resetState()
	input := map[string]any{
		"title": strings.Repeat("t", 500),
	}
	got := Apply(input).(map[string]any)
	if got["title"] != input["title"] {
		t.Error("non-truncatable field should pass through unchanged")
	}
	if _, exists := got["titleLength"]; exists {
		t.Error("non-truncatable field should not get a length annotation")
	}
}

func TestApply_NestedObjects(t *testing.T) {
	resetState()
	long := strings.Repeat("n", 300)
	input := map[string]any{
		"issue": map[string]any{
			"description": long,
			"title":       "Nested Issue",
		},
	}
	got := Apply(input).(map[string]any)
	nested := got["issue"].(map[string]any)
	desc := nested["description"].(string)
	if !strings.HasSuffix(desc, ellipsis) {
		t.Error("expected nested description to be truncated")
	}
}

func TestApply_Arrays(t *testing.T) {
	resetState()
	long := strings.Repeat("a", 300)
	input := []any{
		map[string]any{"description": long},
		map[string]any{"description": "short"},
	}
	got := Apply(input).([]any)
	first := got[0].(map[string]any)
	if !strings.HasSuffix(first["description"].(string), ellipsis) {
		t.Error("expected first item description truncated")
	}
	second := got[1].(map[string]any)
	if second["description"] != "short" {
		t.Error("expected second item description unchanged")
	}
}

func TestApply_Nil(t *testing.T) {
	resetState()
	if Apply(nil) != nil {
		t.Error("Apply(nil) should return nil")
	}
}

func TestApply_NonMapNonSlice(t *testing.T) {
	resetState()
	if Apply(42) != 42 {
		t.Error("Apply on scalar should return unchanged")
	}
	if Apply("hello") != "hello" {
		t.Error("Apply on string should return unchanged")
	}
}

func TestApply_BodyField(t *testing.T) {
	resetState()
	long := strings.Repeat("b", 300)
	input := map[string]any{"body": long}
	got := Apply(input).(map[string]any)
	if !strings.HasSuffix(got["body"].(string), ellipsis) {
		t.Error("expected body field truncated")
	}
	if got["bodyLength"] != 300 {
		t.Errorf("expected bodyLength=300, got %v", got["bodyLength"])
	}
}

func TestApply_ContentField(t *testing.T) {
	resetState()
	long := strings.Repeat("c", 300)
	input := map[string]any{"content": long}
	got := Apply(input).(map[string]any)
	if !strings.HasSuffix(got["content"].(string), ellipsis) {
		t.Error("expected content field truncated")
	}
	if got["contentLength"] != 300 {
		t.Errorf("expected contentLength=300, got %v", got["contentLength"])
	}
}

func TestApply_ExpandCaseInsensitive(t *testing.T) {
	Configure(ConfigOpts{Expand: "DESCRIPTION"})
	defer resetState()

	long := strings.Repeat("d", 500)
	input := map[string]any{"description": long}
	got := Apply(input).(map[string]any)
	if got["description"] != long {
		t.Error("expand should be case-insensitive")
	}
}

func TestApply_ExpandMultipleFields(t *testing.T) {
	Configure(ConfigOpts{Expand: "description,body"})
	defer resetState()

	long := strings.Repeat("e", 500)
	input := map[string]any{
		"description": long,
		"body":        long,
		"content":     long,
	}
	got := Apply(input).(map[string]any)
	if got["description"] != long {
		t.Error("description should be expanded")
	}
	if got["body"] != long {
		t.Error("body should be expanded")
	}
	if !strings.HasSuffix(got["content"].(string), ellipsis) {
		t.Error("content should still be truncated")
	}
}
