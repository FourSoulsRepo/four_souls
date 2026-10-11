package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// countThreePlus counts copies of 3+ player cards in the decks, the
// slots and the discards.
func countThreePlus(g *engine.Game) int {
	var ids []engine.ObjectID
	for d := range g.Decks {
		ids = append(ids, g.Decks[d]...)
		ids = append(ids, g.Discards[d]...)
	}
	for _, s := range append(append([]engine.Slot(nil), g.Shop...), g.Monsters...) {
		ids = append(ids, s.Cards...)
	}
	three := map[engine.CardRef]bool{}
	for _, d := range Set.Cards {
		three[d.Ref] = d.ThreePlus
	}
	n := 0
	for _, id := range ids {
		if three[g.Object(id).Card] {
			n++
		}
	}
	return n
}

// R-CARD-32: cards with the 3+ player symbol are left out of 2-player
// games and in from 3 players.
func TestThreePlusCards(t *testing.T) {
	want := 0
	for _, d := range Set.Cards {
		if d.ThreePlus {
			want += max(d.Copies, 1)
		}
	}
	if want != 11 {
		t.Fatalf("%d 3+ player cards in the Base Game, want 11", want)
	}
	if n := countThreePlus(newGame(t, "isaac", "cain")); n != 0 {
		t.Errorf("a 2-player game has %d 3+ player cards, want none", n)
	}
	if n := countThreePlus(newGame(t, "isaac", "cain", "maggy")); n != want {
		t.Errorf("a 3-player game has %d 3+ player cards, want %d", n, want)
	}
}
