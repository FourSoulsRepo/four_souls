package assets

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// The same file must be readable in both build modes:
// go test ./internal/assets and go test -tags embed ./internal/assets.
func TestAboutFile(t *testing.T) {
	if !Embedded {
		// Tests run in the package folder; point external mode at the source files.
		t.Setenv(EnvDir, "files")
	}
	data, err := fs.ReadFile(FS(), "about.txt")
	if err != nil {
		t.Fatalf("read about.txt (embedded=%v): %v", Embedded, err)
	}
	if len(data) == 0 {
		t.Fatal("about.txt is empty")
	}
}

func TestMissingFolderIsNotFatal(t *testing.T) {
	if Embedded {
		t.Skip("external mode only")
	}
	t.Setenv(EnvDir, filepath.Join(t.TempDir(), "absent"))
	_, err := fs.ReadFile(FS(), "about.txt")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("want not-exist error, got %v", err)
	}
}
