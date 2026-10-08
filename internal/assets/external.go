//go:build !embed

package assets

import "io/fs"

// Embedded reports whether files are packed into the binary.
const Embedded = false

// FS returns the app assets.
func FS() fs.FS {
	return external()
}
