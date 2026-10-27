//go:build embed

package carddb

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed all:data
var embeddedData embed.FS

// Embedded returns the card files packed into the binary (B-01).
// Card data is always packed with -tags embed; images only with
// -tags "embed cardimages", because they come from the private
// card_db submodule that public clones do not have (A-11).
func Embedded() (fs.FS, bool) {
	return union{data: embeddedData, images: embeddedImages()}, true
}

// union serves images/... from the images FS and everything else from data.
// The images FS has no images/ prefix: images/en/b2/x.webp is en/b2/x.webp there.
type union struct {
	data   fs.FS
	images fs.FS // nil when images are not embedded
}

func (u union) Open(name string) (fs.File, error) {
	if name == ImagesDir || strings.HasPrefix(name, ImagesDir+"/") {
		if u.images == nil {
			return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
		}
		sub := strings.TrimPrefix(strings.TrimPrefix(name, ImagesDir), "/")
		if sub == "" {
			sub = "."
		}
		f, err := u.images.Open(sub)
		if err != nil {
			return nil, err //nolint:wrapcheck // keep fs errors as they are
		}
		return f, nil
	}
	f, err := u.data.Open(name)
	if err != nil {
		return nil, err //nolint:wrapcheck // keep fs errors as they are
	}
	return f, nil
}
