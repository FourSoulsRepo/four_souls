//go:build embed && cardimages

package carddb

import (
	"io/fs"

	"github.com/FourSoulsRepo/card_db/images"
)

// The images submodule embeds its own files (8.29).
func embeddedImages() fs.FS {
	return images.FS()
}
