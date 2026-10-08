//go:build !embed

package assets

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Embedded reports whether files are packed into the binary.
const Embedded = false

// dirName is the assets folder name next to the binary.
const dirName = "assets"

// sourceDir is the assets folder in the repo, used by wails dev from the repo root.
const sourceDir = "internal/assets/files"

// FS returns the app assets.
func FS() fs.FS {
	return os.DirFS(externalDir())
}

// externalDir picks the assets folder for external mode.
// Order: EnvDir, an assets folder next to the binary, the repo source folder.
// A missing folder is not an error: every Open then fails with fs.ErrNotExist.
func externalDir() string {
	if dir := os.Getenv(EnvDir); dir != "" {
		return dir
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Join(filepath.Dir(exe), dirName)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	return sourceDir
}
