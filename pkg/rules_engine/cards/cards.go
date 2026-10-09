// Package cards lists every card set the engine ships with.
//
// Each set is its own package with one file per card, e.g. cards/b2.
// Stubs come from card_db: `go run ./cmd/cardgen -set b2`.
package cards

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// Sets returns every known set, in a fixed order.
func Sets() []engine.CardSet {
	return []engine.CardSet{b2.Set}
}

// Find returns the set with the given code, e.g. "b2".
func Find(code string) (engine.CardSet, bool) {
	for _, s := range Sets() {
		if s.Code == code {
			return s, true
		}
	}
	return engine.CardSet{}, false
}
