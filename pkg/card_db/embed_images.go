//go:build embed && cardimages

package carddb

import (
	"embed"
	"io/fs"
)

// The images pattern has no all: prefix, so the submodule's .git is skipped.
//
//go:embed images
var imagesFiles embed.FS

func embeddedImages() fs.FS {
	return imagesFiles
}
