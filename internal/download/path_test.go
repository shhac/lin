package download

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDestPath_DefaultDir(t *testing.T) {
	// DefaultDir (a not-yet-existing cache dir) is created and used when neither
	// --output nor --output-dir is set.
	base := filepath.Join(t.TempDir(), "downloads")
	got, err := resolveDestPath("a.png", DownloadOpts{DefaultDir: base}, "image/png")
	if err != nil {
		t.Fatalf("resolveDestPath: %v", err)
	}
	if want := filepath.Join(base, "a.png"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if _, err := os.Stat(base); err != nil {
		t.Errorf("DefaultDir was not created: %v", err)
	}
}

func TestResolveDestPath_OutputAndDirOverrideDefault(t *testing.T) {
	def := t.TempDir()

	// --output wins.
	out := filepath.Join(t.TempDir(), "explicit.png")
	got, err := resolveDestPath("a.png", DownloadOpts{Output: out, DefaultDir: def}, "image/png")
	if err != nil || got != out {
		t.Errorf("--output: got %q err %v, want %q", got, err, out)
	}

	// --output-dir wins (must exist).
	dir := t.TempDir()
	got, err = resolveDestPath("a.png", DownloadOpts{OutputDir: dir, DefaultDir: def}, "image/png")
	if err != nil || got != filepath.Join(dir, "a.png") {
		t.Errorf("--output-dir: got %q err %v", got, err)
	}
}

func TestResolveDestPath_DefaultDirCreateFails(t *testing.T) {
	// A regular file stands where DefaultDir's parent should be, so MkdirAll
	// can't create the directory — the error must propagate.
	blocker := filepath.Join(t.TempDir(), "afile")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := resolveDestPath("a.png", DownloadOpts{DefaultDir: filepath.Join(blocker, "sub")}, "image/png")
	if err == nil {
		t.Fatal("expected an error when DefaultDir can't be created")
	}
}

func TestResolveDestPath_EmptyDefaultFallsBackToCwd(t *testing.T) {
	got, err := resolveDestPath("a.png", DownloadOpts{}, "image/png")
	if err != nil {
		t.Fatalf("resolveDestPath: %v", err)
	}
	cwd, _ := os.Getwd()
	if want := filepath.Join(cwd, "a.png"); got != want {
		t.Errorf("got %q, want cwd-relative %q", got, want)
	}
}
