package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBlindRage(t *testing.T) {
	tb := roomTable(t, "blind_rage")
	tb.EndTurn()
	if tb.G.Turn.MustAttackDeck != 1 {
		t.Fatal("Cain must attack the monster deck")
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 1, Kind: engine.IntentEndTurn}); err == nil {
		t.Error("Cain can't end the turn before attacking")
	}
}
