package main

import (
	"io/fs"
	"path"

	carddb "github.com/FourSoulsRepo/card_db"
)

// coinIcon is the coin image inside the card images folder: the ¢ face of
// the game's coin, cut out by the tools repo (15.38).
const coinIcon = "assets/icons/coin.webp"

// swordIcon is the game's pixel sword, the Attack icon (15.24).
const swordIcon = "assets/icons/sword.webp"

// bannerImage is the foursouls.ru logo, the home screen banner (15.41).
const bannerImage = "assets/banner.webp"

// Icons are the game's own small images; each is "" when its file is
// missing, and the frontend draws its own.
type Icons struct {
	Coin   string `json:"coin"`   // coin image URL, next to a player's cents
	Sword  string `json:"sword"`  // sword image URL, on the Attack icon
	Banner string `json:"banner"` // banner image URL, atop the home screen
}

// Icons returns the icon URLs.
func (a *App) Icons() Icons {
	return Icons{Coin: a.imageURL(coinIcon), Sword: a.imageURL(swordIcon), Banner: a.imageURL(bannerImage)}
}

// imageURL is the URL of file inside the card images folder, or "" when
// the file is missing.
func (a *App) imageURL(file string) string {
	info, err := fs.Stat(a.cardFiles, path.Join(carddb.ImagesDir, file))
	if err != nil || info.IsDir() {
		return ""
	}
	return cardImagePrefix + file
}
