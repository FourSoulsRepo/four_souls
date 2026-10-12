package main

import (
	"io/fs"
	"path"

	carddb "github.com/FourSoulsRepo/card_db"
)

// DiceSkin is one look of the six-sided die (V-19).
type DiceSkin struct {
	ID   string `json:"id"` // one of settings.Dice
	Name string `json:"name"`
	// Image is the texture URL: a 3×3 grid, top row unused, then faces
	// 2 4 5 and 1 3 6. "" when the file is missing; the die is then
	// drawn plain.
	Image string `json:"image"`
}

// diceSkins are the dice textures inside the card images folder; the
// order and IDs follow settings.Dice.
var diceSkins = []struct{ id, name, file string }{
	{"tear", "Tear", "assets/dices/holy_set/d6_tear.png"},
}

// Dice returns every dice look, for the dice component and settings.
func (a *App) Dice() []DiceSkin {
	out := make([]DiceSkin, 0, len(diceSkins))
	for _, d := range diceSkins {
		skin := DiceSkin{ID: d.id, Name: d.name}
		if info, err := fs.Stat(a.cardFiles, path.Join(carddb.ImagesDir, d.file)); err == nil && !info.IsDir() {
			skin.Image = cardImagePrefix + d.file
		}
		out = append(out, skin)
	}
	return out
}
