package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAmbush(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("ambush")
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err == nil {
		t.Error("ended the turn without the 2 attacks")
	}
	if n := tb.G.Turn.MustAttackDeck; n != 2 {
		t.Errorf("must attack the deck %d times, want 2", n)
	}
}
