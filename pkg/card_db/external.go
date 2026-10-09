//go:build !embed

package carddb

import "io/fs"

// Embedded reports that card files are read from disk in this build (B-01).
func Embedded() (fs.FS, bool) {
	return nil, false
}
