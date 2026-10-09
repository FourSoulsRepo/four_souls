package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// seat is a player with a character and items.
var seat = enginetest.Seat

// table builds a b2 table; player 0 is active.
func table(t *testing.T, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.New(t, Set, seats...)
}

// testExtraLoot checks a character's ↷: one more loot play, on your own
// turn and on another player's turn.
func testExtraLoot(t *testing.T, character engine.CardRef) {
	t.Helper()
	tb := table(t, seat(character), seat(character))
	tb.Start(0, character, 0)
	tb.Pass(0)
	tb.Activate(1, character, 0) // in response, on another player's turn
	if got := tb.G.Turn.LootPlays; got != 2 {
		t.Errorf("active player's loot plays = %d, want 2", got)
	}
	if tb.G.Object(tb.Find(0, character)).Charged {
		t.Error("the character is still charged after ↷")
	}
	if got := tb.G.Players[1].ExtraLootPlays; got != 1 {
		t.Errorf("other player's extra loot plays = %d, want 1", got)
	}
	if rule := tb.Refused(0, character, 0); rule != "R-ABIL-07" {
		t.Errorf("second use refused by %s, want R-ABIL-07 (deactivated)", rule)
	}
}

// newGame starts a real game with fixed characters.
func newGame(t *testing.T, characters ...engine.CardRef) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: len(characters), Sets: []engine.CardSet{Set}, Characters: characters})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// hasItem reports whether player p controls an item with this card.
func hasItem(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].InPlay {
		if o := g.Object(id); o.Card == card && o.Role == engine.RoleItem {
			return true
		}
	}
	return false
}

// monsterTable builds a b2 table with the given monsters in slots.
func monsterTable(t *testing.T, monsters []engine.CardRef, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{Players: seats, Monsters: monsters}, Set)
}
