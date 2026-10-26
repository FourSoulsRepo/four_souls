package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHandMeDowns(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "hand_me_downs"), seat("cain"), seat("samson")}})
	tb.G.PutOnDeck(engine.LootDeck, "bomb")
	tb.Activate(0, "hand_me_downs", 0, "ask.yes", "ask.no", "bomb")
	if len(tb.G.Players[0].Hand) != 1 || len(tb.G.Players[1].Hand) != 0 || len(tb.G.Players[2].Hand) != 0 {
		t.Errorf("hands %d %d %d; Cain loots the bomb and gives it, Samson declines",
			len(tb.G.Players[0].Hand), len(tb.G.Players[1].Hand), len(tb.G.Players[2].Hand))
	}
}
