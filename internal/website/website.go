// Package website is what the rules website's .wasm answers (ADR 026).
// cmd/website exports these functions to JavaScript; here they are
// plain Go, so tests run without a browser.
package website

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"slices"

	carddb "github.com/FourSoulsRepo/card_db"
	"github.com/FourSoulsRepo/four_souls/internal/cardinfo"
	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// Sets are the set codes the site carries: finished sets only, with no
// stub card left. A cmd/cardgen test checks that (ADR 026).
// A set joins when step 14 finishes it.
var Sets = []string{"b2", "g2", "fsp2", "r"}

// Version is the engine and app versions, as JSON.
func Version() string { return encode(version.Get()) }

// SetsJSON is Sets as JSON.
func SetsJSON() string { return encode(Sets) }

// Cards is every card of Sets, as JSON, without images: the site shows
// cards as text. fsys is the card data (card_db's data/ folder).
func Cards(fsys fs.FS) string {
	all, err := cardinfo.List(fsys, nil)
	if err != nil {
		return failure(err)
	}
	out := make([]cardinfo.CardInfo, 0, len(all))
	for _, c := range all {
		if slices.Contains(Sets, c.Set) {
			out = append(out, c)
		}
	}
	return encode(out)
}

// Check answers a situation payload (LR-04) with the site's card sets.
func Check(payload string) string {
	sets := make([]engine.CardSet, 0, len(Sets))
	for _, code := range Sets {
		s, ok := cards.Find(code)
		if !ok {
			return failure(fmt.Errorf("unknown card set %q", code))
		}
		sets = append(sets, s)
	}
	ans, err := engine.RunSituation([]byte(payload), sets...)
	if err != nil {
		return failure(err)
	}
	return encode(ans)
}

// Notice is the fan-game notice (L-01), as JSON.
func Notice() string { return encode(legal.Notice) }

// Recover turns a panic into an error answer, so nothing panics into
// the page. Use it as: defer website.Recover(&out).
func Recover(out *string) {
	if r := recover(); r != nil {
		*out = failure(fmt.Errorf("internal error: %v", r))
	}
}

// failure is an error answer: {"error": "…"}.
func failure(err error) string {
	return encode(struct {
		Error string `json:"error"`
	}{err.Error()})
}

func encode(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return failure(errors.New("encode: " + err.Error()))
	}
	return string(b)
}

// EmbeddedCards is Cards of the card data built into the binary, or an
// error answer when the build has none (built without -tags embed).
func EmbeddedCards() string {
	fsys, ok := carddb.Embedded()
	if !ok {
		return failure(errors.New("no card data: build with -tags embed"))
	}
	return Cards(fsys)
}
