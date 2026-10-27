package assets

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	carddb "github.com/FourSoulsRepo/card_db"
)

// EnvCardsDir overrides the card data folder in external mode.
const EnvCardsDir = "FOUR_SOULS_CARDS"

// EnvImagesDir names a folder of extra card images, such as a
// translator's own repo: <lang>/<set>/<file>, as in images/ (16.12).
// It works in every build mode, embedded too.
const EnvImagesDir = "FOUR_SOULS_IMAGES"

// cardsDirName is the card data folder name next to the binary.
const cardsDirName = "cards"

// cardsSourceDir is the card_db module in the repo, used by wails dev.
const cardsSourceDir = "pkg/card_db"

// CardsFS returns the card data and images: data/*.json and images/.
// With -tags embed they are packed into the binary; otherwise they are read
// from FOUR_SOULS_CARDS, a cards folder next to the binary, or pkg/card_db.
// Missing images are not an error: cards then show as text (A-11).
// FOUR_SOULS_IMAGES lays its images over them.
func CardsFS() fs.FS {
	return withImages(cardsFS(), os.Getenv(EnvImagesDir))
}

// withImages lays the card images in dir over fsys. os.Root keeps every
// read inside dir, symbolic links too; a missing dir is logged and skipped.
func withImages(fsys fs.FS, dir string) fs.FS {
	if dir == "" {
		return fsys
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		log.Printf("card images: %v; %s is skipped", err, EnvImagesDir)
		return fsys
	}
	return carddb.WithImages(fsys, root.FS())
}

func cardsFS() fs.FS {
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
