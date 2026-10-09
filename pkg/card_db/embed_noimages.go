//go:build embed && !cardimages

package carddb

import "io/fs"

// embeddedImages is nil without -tags cardimages: cards show as text.
func embeddedImages() fs.FS {
	return nil
}
