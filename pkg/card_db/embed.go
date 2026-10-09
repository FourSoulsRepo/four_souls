//go:build embed

package carddb

import (
	"embed"
	"io/fs"
)

//go:embed all:data all:images
var embedded embed.FS

// Embedded returns the card files packed into the binary (B-01).
func Embedded() (fs.FS, bool) {
	return embedded, true
}
