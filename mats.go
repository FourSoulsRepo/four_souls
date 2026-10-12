package main

import (
	"io/fs"
	"path"
	"strings"

	carddb "github.com/FourSoulsRepo/card_db"

	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

// matsDir holds the game mat images inside the card images folder.
const matsDir = "assets/mats"

// MatInfo is one game mat image (ST-03, V-24).
type MatInfo struct {
	ID    string `json:"id"` // the file name without ".webp"; the Mat setting
	Name  string `json:"name"`
	Image string `json:"image"`
}

// Mats lists every mat image, so a new mat needs only a new file.
// Without the private images the list is empty and the table shows
// the color mats.
func (a *App) Mats() []MatInfo {
	out := []MatInfo{}
	entries, err := fs.ReadDir(a.cardFiles, path.Join(carddb.ImagesDir, matsDir))
	if err != nil {
		return out
	}
	for _, e := range entries {
		id, ok := strings.CutSuffix(e.Name(), ".webp")
		if e.IsDir() || !ok || !settings.ValidMat(id) {
			continue
		}
		out = append(out, MatInfo{ID: id, Name: matTitle(id), Image: cardImagePrefix + matsDir + "/" + e.Name()})
	}
	return out
}

// matTitle turns a file name into a menu name: "dark_room" is "Dark room".
func matTitle(id string) string {
	name := strings.NewReplacer("_", " ", "-", " ").Replace(id)
	return strings.ToUpper(name[:1]) + name[1:]
}
