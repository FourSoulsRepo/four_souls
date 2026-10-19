package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// Requiem reprints A Sack, a Gold Box card. A game without Gold Box
// borrows its definition and has only Requiem's copy (ADR 010).
func TestReprintWithoutItsSet(t *testing.T) {
	g, _, err := engine.NewGame(engine.Setup{Seed: 3, Players: 2, Sets: []engine.CardSet{b2.Set, Set}})
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, o := range g.Objects {
		if o.Card == "a_sack" {
			n++
		}
	}
	if n != 1 {
		t.Errorf("A Sack copies = %d, want Requiem's 1", n)
	}
}
