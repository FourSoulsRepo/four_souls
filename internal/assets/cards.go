package assets

import (
	"io/fs"
	"os"
	"path/filepath"

	carddb "github.com/FourSoulsRepo/card_db"
)

// EnvCardsDir overrides the card data folder in external mode.
const EnvCardsDir = "FOUR_SOULS_CARDS"

// cardsDirName is the card data folder name next to the binary.
const cardsDirName = "cards"

// cardsSourceDir is the card_db module in the repo, used by wails dev.
const cardsSourceDir = "pkg/card_db"

// CardsFS returns the card data and images: data/*.json and images/.
// With -tags embed they are packed into the binary; otherwise they are read
// from FOUR_SOULS_CARDS, a cards folder next to the binary, or pkg/card_db.
// Missing images are not an error: cards then show as text (A-11).
func CardsFS() fs.FS {
	if fsys, ok := carddb.Embedded(); ok {
		return fsys
	}
	if dir := os.Getenv(EnvCardsDir); dir != "" {
		return os.DirFS(dir)
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Join(filepath.Dir(exe), cardsDirName)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return os.DirFS(dir)
		}
	}
	return os.DirFS(cardsSourceDir)
}
