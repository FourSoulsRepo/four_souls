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

// skyLayers are the foursouls.ru page background (15.42): the sky, then
// its see-through star layers, far to near; each a 600×400 tile.
var skyLayers = []string{
	"assets/background/sky.webp",
	"assets/background/stars-1.webp",
	"assets/background/stars-2.webp",
	"assets/background/stars-3.webp",
	"assets/background/stars-4.webp",
	"assets/background/stars-5.webp",
}

// Icons are the game's own small images; each is "" when its file is
// missing, and the frontend draws its own.
type Icons struct {
	Coin   string `json:"coin"`   // coin image URL, next to a player's cents
	Sword  string `json:"sword"`  // sword image URL, on the Attack icon
	Banner string `json:"banner"` // banner image URL, atop the home screen
	// Sky is the menu background: the URLs of skyLayers, or empty when
	// any is missing, and the menus keep their plain color.
	Sky []string `json:"sky"`
}

// Icons returns the icon URLs.
func (a *App) Icons() Icons {
	return Icons{Coin: a.imageURL(coinIcon), Sword: a.imageURL(swordIcon), Banner: a.imageURL(bannerImage), Sky: a.sky()}
}

// sky returns the URLs of every sky layer, or an empty list when one is
// missing: a sky without some of its stars would not match the site.
func (a *App) sky() []string {
	urls := make([]string, 0, len(skyLayers))
	for _, file := range skyLayers {
		u := a.imageURL(file)
		if u == "" {
			return []string{}
		}
		urls = append(urls, u)
	}
	return urls
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
