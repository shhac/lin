package upload

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPPutWithHeaders_OK(t *testing.T) {
	var seenMethod string
	var seenLength int64
	var seenHeader string
	var seenBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenMethod = r.Method
		seenLength = r.ContentLength
		seenHeader = r.Header.Get("X-Test")
		body, _ := io.ReadAll(r.Body)
		seenBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	body := strings.NewReader("payload-bytes")
	err := httpPutWithHeaders(srv.URL, body, int64(body.Len()), map[string]string{"X-Test": "yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seenMethod != "PUT" {
		t.Errorf("method = %q, want PUT", seenMethod)
	}
	if seenLength != 13 {
		t.Errorf("ContentLength = %d, want 13", seenLength)
	}
	if seenHeader != "yes" {
		t.Errorf("X-Test header = %q, want %q", seenHeader, "yes")
	}
	if seenBody != "payload-bytes" {
		t.Errorf("body = %q", seenBody)
	}
}

func TestHTTPPutWithHeaders_RemoteError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	body := strings.NewReader("payload")
	err := httpPutWithHeaders(srv.URL, body, int64(body.Len()), nil)
	if err == nil {
		t.Fatal("expected error for 403")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error = %v, want to contain 403", err)
	}
}

func TestFormatFileMarkdown_ImageTypes(t *testing.T) {
	files := []UploadedFile{
		{Filename: "screenshot.png", AssetURL: "https://example.com/a.png", ContentType: "image/png"},
		{Filename: "photo.jpg", AssetURL: "https://example.com/b.jpg", ContentType: "image/jpeg"},
	}
	got := FormatFileMarkdown(files)
	want := "![screenshot.png](https://example.com/a.png)\n![photo.jpg](https://example.com/b.jpg)"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatFileMarkdown_NonImageTypes(t *testing.T) {
	files := []UploadedFile{
		{Filename: "report.pdf", AssetURL: "https://example.com/r.pdf", ContentType: "application/pdf"},
		{Filename: "data.csv", AssetURL: "https://example.com/d.csv", ContentType: "text/csv"},
	}
	got := FormatFileMarkdown(files)
	want := "[report.pdf](https://example.com/r.pdf)\n[data.csv](https://example.com/d.csv)"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatFileMarkdown_Mixed(t *testing.T) {
	files := []UploadedFile{
		{Filename: "diagram.png", AssetURL: "https://example.com/d.png", ContentType: "image/png"},
		{Filename: "notes.txt", AssetURL: "https://example.com/n.txt", ContentType: "text/plain"},
	}
	got := FormatFileMarkdown(files)
	want := "![diagram.png](https://example.com/d.png)\n[notes.txt](https://example.com/n.txt)"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatFileMarkdown_Empty(t *testing.T) {
	got := FormatFileMarkdown(nil)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestDetectMIME_CommonExtensions(t *testing.T) {
	tests := []struct {
		filename string
		wantType string
	}{
		{"photo.png", "image/png"},
		{"doc.pdf", "application/pdf"},
		{"data.json", "application/json"},
		{"style.css", "text/css; charset=utf-8"},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := detectMIME(tt.filename)
			if got != tt.wantType {
				t.Errorf("detectMIME(%q) = %q, want %q", tt.filename, got, tt.wantType)
			}
		})
	}
}

func TestDetectMIME_UnknownExtension(t *testing.T) {
	got := detectMIME("file.xyz123")
	if got != "application/octet-stream" {
		t.Errorf("expected fallback content type, got %q", got)
	}
}

func TestDetectMIME_NoExtension(t *testing.T) {
	got := detectMIME("Makefile")
	if got != "application/octet-stream" {
		t.Errorf("expected fallback for no extension, got %q", got)
	}
}
