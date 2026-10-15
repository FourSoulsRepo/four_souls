package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"sync"

	carddb "github.com/FourSoulsRepo/card_db"
)

// cardImagePrefix is the URL path card images are served under.
const cardImagePrefix = "/cardimg/"

// CardInfo is what the frontend shows of one card (7.4).
type CardInfo struct {
	Ref       string            `json:"ref"` // as in the engine's views, e.g. "the_d6"
	Set       string            `json:"set"`
	Type      string            `json:"type"` // e.g. "Eternal Treasure Card"
	Name      string            `json:"name"`
	Effects   []string          `json:"effects"`
	Footnotes []carddb.Footnote `json:"footnotes"`
	Stats     []carddb.Stat     `json:"stats"`
	Rewards   []carddb.Reward   `json:"rewards"`
	// Image is the image URL; "" when the file is missing, and the card
	// shows as a blank card with text (A-11).
	Image string `json:"image"`
}

// cardList loads the card data once.
type cardList struct {
	once  sync.Once
	cards []CardInfo
	err   error
}

func (l *cardList) load(fsys fs.FS) ([]CardInfo, error) {
	l.once.Do(func() {
		db, err := carddb.Load(fsys)
		if err != nil {
			l.err = fmt.Errorf("card data: %w", err)
			return
		}
		for _, ref := range db.Refs() {
			c, _ := db.Get(ref)
			en := c.Text[carddb.English]
			// Lists are never nil, so the frontend never gets null.
			info := CardInfo{
				Ref: ref.String(), Set: c.Set, Type: c.Type, Name: en.Name,
				Effects:   append([]string{}, en.Effects...),
				Footnotes: append([]carddb.Footnote{}, en.Footnotes...),
				Stats:     append([]carddb.Stat{}, c.Stats...),
				Rewards:   append([]carddb.Reward{}, c.Rewards...),
			}
			if carddb.HasImage(fsys, c) {
				info.Image = cardImagePrefix + c.Images[0]
			}
			l.cards = append(l.cards, info)
		}
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
