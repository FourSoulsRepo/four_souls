// Package assets is the binding of the card list and the game's media:
// card backs, dice, mats, fonts, icons and sounds (ADR 030). It also
// serves their files to the frontend.
package assets

import (
	"io/fs"
	"path"

	carddb "github.com/FourSoulsRepo/card_db"
	"github.com/FourSoulsRepo/four_souls/internal/cardinfo"
)

// cardImagePrefix is the URL path card images are served under.
const cardImagePrefix = cardinfo.ImagePrefix

// Assets tells the frontend which cards and media files the app has.
type Assets struct {
	cardFiles fs.FS // card data and images
	cardList  cardList
}

// New makes the binding over the card files: data/*.json and images/.
func New(cardFiles fs.FS) *Assets {
	return &Assets{cardFiles: cardFiles}
}

// imageURL is the URL of file inside the card images folder, or "" when
// the file is missing.
func (a *Assets) imageURL(file string) string {
	info, err := fs.Stat(a.cardFiles, path.Join(carddb.ImagesDir, file))
	if err != nil || info.IsDir() {
		return ""
	}
	return cardImagePrefix + file
}
