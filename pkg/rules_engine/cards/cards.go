// Package cards lists every card set the engine ships with.
//
// Each set is its own package with one file per card, e.g. cards/b2.
// Stubs come from card_db: `go run ./cmd/cardgen -set b2`.
package cards

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
	"github.com/FourSoulsRepo/rules_engine/cards/fsp2"
	"github.com/FourSoulsRepo/rules_engine/cards/g2"
	"github.com/FourSoulsRepo/rules_engine/cards/r"
)

// Sets returns every known set, in a fixed order.
func Sets() []engine.CardSet {
	return []engine.CardSet{b2.Set, g2.Set, fsp2.Set, r.Set}
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

// Challenges returns every known challenge, in a fixed order (ADR 028).
// Each one's cards are in the set its Set names, which Find does not
// list: the caller adds it to the game's sets.
func Challenges() []*engine.ChallengeDef {
	return []*engine.ChallengeDef{&challenge.Krampus}
}

// FindChallenge returns the challenge with the given ID, e.g. "krampus".
func FindChallenge(id string) (*engine.ChallengeDef, bool) {
	for _, c := range Challenges() {
		if c.ID == id {
			return c, true
		}
	}
	return nil, false
}

// ChallengeSet returns the card set that holds a challenge's cards
// (ChallengeDef.Set). Find leaves it out: it is never a playable set.
func ChallengeSet(c *engine.ChallengeDef) (engine.CardSet, bool) {
	for _, s := range []engine.CardSet{challenge.Set} {
		if s.Code == c.Set {
			return s, true
		}
	}
	return engine.CardSet{}, false
}
