package download

import (
	"strings"
	"testing"
)

func TestParseFileURL_FullHTTPS(t *testing.T) {
	input := "https://uploads.linear.app/a1b2c3d4-e5f6-7890-abcd-ef1234567890/11111111-2222-3333-4444-555555555555/66666666-7777-8888-9999-aaaaaaaaaaaa"
	got, err := ParseFileURL(input, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OrgID != "a1b2c3d4-e5f6-7890-abcd-ef1234567890" {
		t.Errorf("OrgID = %q", got.OrgID)
	}
	if len(got.Segments) != 3 {
		t.Errorf("expected 3 segments, got %d", len(got.Segments))
	}
	if !strings.HasPrefix(got.URL, "https://uploads.linear.app/") {
		t.Errorf("URL = %q", got.URL)
	}
}

func TestParseFileURL_HostRelative(t *testing.T) {
	input := "uploads.linear.app/a1b2c3d4-e5f6-7890-abcd-ef1234567890/11111111-2222-3333-4444-555555555555/66666666-7777-8888-9999-aaaaaaaaaaaa"
	got, err := ParseFileURL(input, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OrgID != "a1b2c3d4-e5f6-7890-abcd-ef1234567890" {
		t.Errorf("OrgID = %q", got.OrgID)
	}
}

func TestParseFileURL_PathOnly(t *testing.T) {
	orgID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	input := "11111111-2222-3333-4444-555555555555/66666666-7777-8888-9999-aaaaaaaaaaaa"
	got, err := ParseFileURL(input, orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OrgID != orgID {
		t.Errorf("expected default OrgID, got %q", got.OrgID)
	}
	if len(got.Segments) != 3 {
		t.Errorf("expected 3 segments, got %d", len(got.Segments))
	}
}

func TestParseFileURL_SingleUUID(t *testing.T) {
	orgID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	input := "11111111-2222-3333-4444-555555555555"
	got, err := ParseFileURL(input, orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Segments) != 2 {
		t.Errorf("expected 2 segments (org + file), got %d", len(got.Segments))
	}
}

func TestParseFileURL_HTTP_Rejected(t *testing.T) {
	input := "http://uploads.linear.app/a1b2c3d4-e5f6-7890-abcd-ef1234567890/11111111-2222-3333-4444-555555555555"
	_, err := ParseFileURL(input, "")
	if err == nil {
		t.Fatal("expected error for http://")
	}
	if !strings.Contains(err.Error(), "http://") {
		t.Errorf("error should mention http://, got: %v", err)
	}
}

func TestParseFileURL_WrongHost(t *testing.T) {
	input := "https://evil.example.com/a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	_, err := ParseFileURL(input, "")
	if err == nil {
		t.Fatal("expected error for wrong host")
	}
	if !strings.Contains(err.Error(), "Invalid host") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseFileURL_NonUUIDSegments(t *testing.T) {
	input := "not-a-uuid/also-not-uuid"
	_, err := ParseFileURL(input, "a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	if err == nil {
		t.Fatal("expected error for non-UUID segments")
	}
	if !strings.Contains(err.Error(), "Invalid UUID segment") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseFileURL_TooManySegments(t *testing.T) {
	input := "a1b2c3d4-e5f6-7890-abcd-ef1234567890/11111111-2222-3333-4444-555555555555/22222222-3333-4444-5555-666666666666/33333333-4444-5555-6666-777777777777"
	_, err := ParseFileURL(input, "")
	if err == nil {
		t.Fatal("expected error for too many segments")
	}
	if !strings.Contains(err.Error(), "1-3 UUID") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseFileURL_NoOrgID_NoDefault(t *testing.T) {
	input := "11111111-2222-3333-4444-555555555555"
	_, err := ParseFileURL(input, "")
	if err == nil {
		t.Fatal("expected error when no org ID and no default")
	}
	if !strings.Contains(err.Error(), "organization ID") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSanitizeFilename_Normal(t *testing.T) {
	got := SanitizeFilename("report.pdf")
	if got != "report.pdf" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestSanitizeFilename_PathSeparators(t *testing.T) {
	got := SanitizeFilename("/path/to/file.txt")
	if got != "file.txt" {
		t.Errorf("expected stripped path, got %q", got)
	}
	got = SanitizeFilename("C:\\Users\\doc.pdf")
	if got != "doc.pdf" {
		t.Errorf("expected stripped Windows path, got %q", got)
	}
}

func TestSanitizeFilename_UnsafeChars(t *testing.T) {
	got := SanitizeFilename("file<name>.txt")
	if strings.ContainsAny(got, "<>") {
		t.Errorf("unsafe chars should be replaced, got %q", got)
	}
}

func TestSanitizeFilename_Empty(t *testing.T) {
	got := SanitizeFilename("")
	if got != "download" {
		t.Errorf("expected 'download' for empty, got %q", got)
	}
}

func TestSanitizeFilename_LongName(t *testing.T) {
	long := strings.Repeat("a", 300)
	got := SanitizeFilename(long)
	if len(got) > 255 {
		t.Errorf("expected max 255 chars, got %d", len(got))
	}
}

func TestParseContentDispositionFilename_RFC5987(t *testing.T) {
	header := `attachment; filename*=UTF-8''my%20document.pdf`
	got := parseContentDispositionFilename(header)
	if got != "my document.pdf" {
		t.Errorf("expected decoded filename, got %q", got)
	}
}

func TestParseContentDispositionFilename_Quoted(t *testing.T) {
	header := `attachment; filename="report-final.xlsx"`
	got := parseContentDispositionFilename(header)
	if got != "report-final.xlsx" {
		t.Errorf("expected quoted filename, got %q", got)
	}
}

func TestParseContentDispositionFilename_Unquoted(t *testing.T) {
	header := `attachment; filename=simple.txt`
	got := parseContentDispositionFilename(header)
	if got != "simple.txt" {
		t.Errorf("expected unquoted filename, got %q", got)
	}
}

func TestParseContentDispositionFilename_NoFilename(t *testing.T) {
	header := `inline`
	got := parseContentDispositionFilename(header)
	if got != "" {
		t.Errorf("expected empty for no filename, got %q", got)
	}
}

func TestParseContentDispositionFilename_Empty(t *testing.T) {
	got := parseContentDispositionFilename("")
	if got != "" {
		t.Errorf("expected empty for empty header, got %q", got)
	}
}

func TestMimeToExtension(t *testing.T) {
	tests := []struct {
		mime string
		want string
	}{
		{"image/png", ".png"},
		{"image/jpeg", ".jpg"},
		{"application/pdf", ".pdf"},
		{"text/plain", ".txt"},
		{"application/json", ".json"},
		{"image/png; charset=utf-8", ".png"},
		{"unknown/type", ""},
	}
	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			got := mimeToExtension(tt.mime)
			if got != tt.want {
				t.Errorf("mimeToExtension(%q) = %q, want %q", tt.mime, got, tt.want)
			}
		})
	}
}
