package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestKrampus(t *testing.T) {
	tb := fight(t, "krampus", "fly")
	tb.EndTurn()
	if _, err := tb.G.Apply(engine.Intent{Player: tb.G.Turn.Active, Kind: engine.IntentEndTurn}); err == nil {
		t.Error("Cain must attack Krampus before ending the turn")
	}
}

// Revealed mid-turn by a refill, Krampus must be attacked that turn
// (19.25).
func TestKrampusRevealedMidTurn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.Turn.Attacks = 2
	tb.G.PutOnDeck(engine.MonsterDeck, "krampus")
	tb.Attack("fly", 6) // the fly dies; Krampus fills its slot
	if top(tb, 0) != "krampus" {
		t.Fatalf("slot 1 = %s, want krampus", top(tb, 0))
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err == nil {
		t.Error("with an attack left, Isaac must attack Krampus before ending the turn")
	}
}
