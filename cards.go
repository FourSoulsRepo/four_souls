package main

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
	"sync"

	carddb "github.com/FourSoulsRepo/card_db"
	"github.com/FourSoulsRepo/four_souls/internal/cardinfo"
)

// cardImagePrefix is the URL path card images are served under.
const cardImagePrefix = cardinfo.ImagePrefix

// CardInfo is what the frontend shows of one card (7.4).
type CardInfo = cardinfo.CardInfo

// cardList loads the card data once.
type cardList struct {
	once  sync.Once
	cards []CardInfo
	err   error
}

func (l *cardList) load(fsys fs.FS) ([]CardInfo, error) {
	l.once.Do(func() {
		l.cards, l.err = cardinfo.List(fsys, carddb.ImageLangs(fsys))
	})
	return l.cards, l.err
}

// Cards returns every card the app knows, for the card component.
func (a *App) Cards() ([]CardInfo, error) {
	return a.cardList.load(a.cardFiles)
}

// imageTypes are the files cardImages serves, with their content types:
// card images, card backs, dice textures (V-19) and mats are all WebP;
// the game's fonts are TrueType (15.20); the damage sounds are MP3
// (15.50). The type is set here, as not every OS's MIME table knows .ttf.
var imageTypes = map[string]string{".webp": "image/webp", ".ttf": "font/ttf", ".mp3": "audio/mpeg"}

// cardImages serves card images, backs, dice textures, mats, fonts and sounds under
// cardImagePrefix from the card files; every other request goes on to the frontend.
func cardImages(fsys fs.FS) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name, ok := strings.CutPrefix(r.URL.Path, cardImagePrefix)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			kind, known := imageTypes[path.Ext(name)]
			if !fs.ValidPath(name) || !known {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", kind)
			http.ServeFileFS(w, r, fsys, path.Join(carddb.ImagesDir, name)) //nolint:gosec // fs.ValidPath above rejects ".." and absolute paths
		})
	}
}
