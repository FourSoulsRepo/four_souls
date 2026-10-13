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
