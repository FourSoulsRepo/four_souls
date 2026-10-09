package carddb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Lang is a text language code, e.g. "en".
type Lang string

// English is the only language of the first version (LO-01).
const English Lang = "en"

// Set describes one card set file.
type Set struct {
	Code string `json:"code"` // e.g. "b2"
	Name string `json:"name"` // e.g. "Base Game V2"
}

// Card is the display data of one card version.
type Card struct {
	ID      string `json:"id"`                // stable, e.g. "the_d6"
	Version int    `json:"version,omitempty"` // 0 or 1 = first version (CD-08)
	Set     string `json:"-"`                 // set code, filled by the loader
	Type    string `json:"type"`              // site label, e.g. "Eternal Treasure Card"
	Copies  int    `json:"copies"`            // copies in the set

	Text map[Lang]Text `json:"text"`

	Stats   []Stat   `json:"stats,omitempty"`
	Rewards []Reward `json:"rewards,omitempty"`

	Artists     []Artist          `json:"artists,omitempty"`     // names only (CD-04)
	Translators map[Lang][]string `json:"translators,omitempty"` // names only (CD-06)

	Image      string   `json:"image,omitempty"` // path inside the images folder
	Related    []string `json:"related,omitempty"`
	Rebalanced bool     `json:"rebalanced,omitempty"`
}

// Text is everything printed on a card in one language.
// Icons inside text are written as {Icon Name}, e.g. {Tap Effect}.
type Text struct {
	Name      string     `json:"name"`
	Effects   []string   `json:"effects,omitempty"`
	Footnotes []Footnote `json:"footnotes,omitempty"`
	Notes     []string   `json:"notes,omitempty"`
}

// Footnote is a keyword box, e.g. "-Eternal-".
type Footnote struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Stat is a printed stat, e.g. HP 2. Values stay strings ("X" exists).
type Stat struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Reward is a printed reward, e.g. 1x Loot.
type Reward struct {
	Count string `json:"count"`
	Kind  string `json:"kind"`
}

// Artist credits a person by name.
type Artist struct {
	Role string `json:"role"`
	Name string `json:"name"`
}

// Ref names one card version: "the_d6" or "the_d6@2".
type Ref struct {
	ID      string
	Version int // 1 for the first version
}

var idPattern = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)

// ErrBadRef means a reference does not have the form id or id@N.
var ErrBadRef = errors.New("carddb: bad card reference")

// ParseRef reads "id" or "id@N" with N >= 2.
func ParseRef(s string) (Ref, error) {
	id, ver, hasVer := strings.Cut(s, "@")
	if !idPattern.MatchString(id) {
		return Ref{}, fmt.Errorf("%w: %q", ErrBadRef, s)
	}
	if !hasVer {
		return Ref{ID: id, Version: 1}, nil
	}
	n, err := strconv.Atoi(ver)
	if err != nil || n < 2 || ver != strconv.Itoa(n) {
		return Ref{}, fmt.Errorf("%w: %q", ErrBadRef, s)
	}
	return Ref{ID: id, Version: n}, nil
}

// String formats the reference; the first version has no suffix.
func (r Ref) String() string {
	if r.Version <= 1 {
		return r.ID
	}
	return r.ID + "@" + strconv.Itoa(r.Version)
}

// Ref returns the reference of the card.
func (c Card) Ref() Ref {
	v := c.Version
	if v < 1 {
		v = 1
	}
	return Ref{ID: c.ID, Version: v}
}

// Name returns the card name in lang, falling back to English.
func (c Card) Name(lang Lang) string {
	if t, ok := c.Text[lang]; ok && t.Name != "" {
		return t.Name
	}
	return c.Text[English].Name
}
