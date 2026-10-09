package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRagMan(t *testing.T) {
	tb := slayTable(t, "rag_man", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "rag_man")
	// On top of the monster deck, then drawn back into the empty slot.
	if top, _ := tb.G.Monsters[0].TopOf(); tb.G.Object(top).Card != "rag_man" {
		t.Error("Rag Man did not come back from the top of the monster deck")
	}
	if s := tb.G.SoulValue(0); s != 0 {
		t.Errorf("soul value %d, want 0: it went back to the deck", s)
	}
	if h := len(tb.G.Players[0].Hand); h != 3 {
		t.Errorf("hand %d, want 3: rewards come first", h)
	}
}
