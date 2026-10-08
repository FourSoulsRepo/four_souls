//go:build embed

package assets

import (
	"embed"
	"io/fs"
)

//go:embed all:files
var files embed.FS

// Embedded reports whether files are packed into the binary.
const Embedded = true

// FS returns the app assets.
func FS() fs.FS {
	sub, err := fs.Sub(files, "files")
	if err != nil {
		panic(err) // "files" is a constant path inside the embed; this cannot fail
	}
	return sub
}
