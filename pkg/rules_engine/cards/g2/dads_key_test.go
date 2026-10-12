package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDadsKey(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"dads_key"}, Hand: []engine.CardRef{"bomb"}}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "charger")
	tb.Activate(0, "dads_key", 0, "bomb", "slot 2")
	if top, _ := tb.G.Monsters[1].TopOf(); tb.G.Object(top).Card != "charger" {
		t.Errorf("slot 2 shows %s, want the charger", tb.G.Object(top).Card)
	}
	if len(tb.G.Players[0].Hand) != 0 {
		t.Error("the bomb was not discarded")
	}
}
