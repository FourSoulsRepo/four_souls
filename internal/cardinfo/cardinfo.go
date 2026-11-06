// Package cardinfo builds what a client shows of each card from card_db
// (7.4). The app and the rules website share it (ADR 026).
package cardinfo

import (
	"fmt"
	"io/fs"

	carddb "github.com/FourSoulsRepo/card_db"
)

// ImagePrefix is the URL path the app serves card images under.
const ImagePrefix = "/cardimg/"

// CardInfo is what the frontend shows of one card (7.4).
type CardInfo struct {
	Ref  string `json:"ref"` // as in the engine's views, e.g. "the_d6"
	Set  string `json:"set"`
	Type string `json:"type"` // e.g. "Eternal Treasure Card"
	Name string `json:"name"` // English; the client shows Names in its language
	// Names holds the name in every language the card data has, English
	// too (LO-02); Translators the names credited per language (CD-06).
	Names       map[string]string   `json:"names"`
	Translators map[string][]string `json:"translators"`
	Effects     []string            `json:"effects"`
	Footnotes   []carddb.Footnote   `json:"footnotes"`
	Stats       []carddb.Stat       `json:"stats"`
	Rewards     []carddb.Reward     `json:"rewards"`
	// Image is the English image URL; "" when the file is missing, and
	// the card shows as a blank card with text (A-11).
	Image string `json:"image"`
	// Images holds the image URL per language that has an image folder,
	// English when the card has no image of its own there (16.6).
	Images map[string]string `json:"images"`
	// Back is the ref of a two-sided card's back face; Front, on a back
	// face, the ref of its front (ADR 023).
	Back  string `json:"back,omitempty"`
	Front string `json:"front,omitempty"`
}

// List loads the card data in fsys and returns every card and back
// face. Images are looked up for langs only; with none, every card
// shows as text (the rules website, ADR 026).
func List(fsys fs.FS, langs []carddb.Lang) ([]CardInfo, error) {
	db, err := carddb.Load(fsys)
	if err != nil {
		return nil, fmt.Errorf("card data: %w", err)
	}
	var cards []CardInfo
	for _, ref := range db.Refs() {
		c, _ := db.Get(ref)
		info := build(fsys, langs, c)
		// A back face is a card of its own, so a flipped card's ref
		// finds it (ADR 023).
		b, two := c.BackCard()
		if two {
			info.Back = b.ID
		}
		cards = append(cards, info)
		if two {
			cards = append(cards, build(fsys, langs, b))
		}
	}
	return cards, nil
}

// build is what the frontend shows of one card or back face.
func build(fsys fs.FS, langs []carddb.Lang, c carddb.Card) CardInfo {
	en := c.Text[carddb.English]
	// Lists are never nil, so the frontend never gets null.
	info := CardInfo{
		Ref: c.Ref().String(), Set: c.Set, Type: c.Type, Name: en.Name, Front: c.Front,
		Effects:   append([]string{}, en.Effects...),
		Footnotes: append([]carddb.Footnote{}, en.Footnotes...),
		Stats:     append([]carddb.Stat{}, c.Stats...),
		Rewards:   append([]carddb.Reward{}, c.Rewards...),
		Names:     map[string]string{}, Translators: map[string][]string{}, Images: map[string]string{},
	}
	for lang, t := range c.Text {
		if t.Name != "" {
			info.Names[string(lang)] = t.Name
		}
	}
	for lang, names := range c.Translators {
		info.Translators[string(lang)] = append([]string{}, names...)
	}
	for _, lang := range langs {
		if p := carddb.FindImage(fsys, c, lang); p != "" {
			info.Images[string(lang)] = ImagePrefix + p
		}
	}
	info.Image = info.Images[string(carddb.English)]
	return info
}
