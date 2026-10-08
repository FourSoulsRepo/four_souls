// Package assets gives access to app files (B-01).
//
// Built with -tags embed, files are packed into the binary.
// Built without it (the default, used by wails dev), files are read
// from an assets folder on disk, so they can change without a rebuild.
package assets

import (
	"io/fs"
	"os"
	"path/filepath"
)

// EnvDir overrides the assets folder in external mode.
const EnvDir = "FOUR_SOULS_ASSETS"

// dirName is the assets folder name next to the binary.
const dirName = "assets"

// sourceDir is the assets folder in the repo, used by wails dev from the repo root.
const sourceDir = "internal/assets/files"

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

// external returns the on-disk assets.
func external() fs.FS {
	return os.DirFS(externalDir())
}
