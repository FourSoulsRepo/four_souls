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
	Set  string `json:"set"` // "Holy" or "Unholy", for the settings menu
	// Image is the texture URL: a 3×3 grid, top row unused, then faces
	// 2 4 5 and 1 3 6. "" when the file is missing; the die is then
	// drawn plain.
	Image string `json:"image"`
}

// diceSkins are the dice textures inside the card images folder; the
// order and IDs follow settings.Dice.
var diceSkins = []struct{ id, name, set, file string }{
	{"tear", "Tear", "Holy", "assets/dices/holy_set/d6_tear.webp"},
	{"blue_fly", "Blue Fly", "Holy", "assets/dices/holy_set/d6_blue_fly.webp"},
	{"cent", "Cent", "Holy", "assets/dices/holy_set/d6_cent.webp"},
	{"moms_knife", "Mom's Knife", "Holy", "assets/dices/holy_set/d6_moms_knife.webp"},
	{"the_bible", "The Bible", "Holy", "assets/dices/holy_set/d6_the_bible.webp"},
	{"brimstone", "Brimstone", "Unholy", "assets/dices/unholy_set/d6_brimstone.webp"},
	{"guppy_paw", "Guppy's Paw", "Unholy", "assets/dices/unholy_set/d6_guppy_paw.webp"},
	{"pentagram", "Pentagram", "Unholy", "assets/dices/unholy_set/d6_pentagram.webp"},
	{"poop", "Poop", "Unholy", "assets/dices/unholy_set/d6_poop.webp"},
	{"tinted_rock", "Tinted Rock", "Unholy", "assets/dices/unholy_set/d6_tinted_rock.webp"},
}

// Dice returns every dice look, for the dice component and settings.
func (a *App) Dice() []DiceSkin {
	out := make([]DiceSkin, 0, len(diceSkins))
	for _, d := range diceSkins {
		skin := DiceSkin{ID: d.id, Name: d.name, Set: d.set}
		if info, err := fs.Stat(a.cardFiles, path.Join(carddb.ImagesDir, d.file)); err == nil && !info.IsDir() {
			skin.Image = cardImagePrefix + d.file
		}
		out = append(out, skin)
	}
	return out
}
