package main

import (
	"io/fs"
	"path"

	carddb "github.com/FourSoulsRepo/card_db"
)

// backsDir holds the card back images inside the card images folder.
const backsDir = "assets/backs"

// backKinds are the decks with a back image, as the frontend's card kinds.
var backKinds = []string{"loot", "treasure", "monster"}

// CardBack is the back image of one deck (V-26).
type CardBack struct {
	Kind string `json:"kind"` // "loot", "treasure" or "monster"
	// Image is the image URL; "" when the file is missing, and the back
	// is drawn plain.
	Image string `json:"image"`
}

// Backs returns the back of every deck, for face-down cards.
func (a *App) Backs() []CardBack {
	out := make([]CardBack, 0, len(backKinds))
	for _, kind := range backKinds {
		b := CardBack{Kind: kind}
		file := backsDir + "/" + kind + ".webp"
		if info, err := fs.Stat(a.cardFiles, path.Join(carddb.ImagesDir, file)); err == nil && !info.IsDir() {
			b.Image = cardImagePrefix + file
		}
		out = append(out, b)
	}
	return out
}
