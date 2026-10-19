package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBookOfVirtues(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("bethany", "book_of_virtues"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.ForceRolls(6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "leech", foe, foe)
	if top(tb, 1) == "leech" {
		t.Error("its 5-6: 1 damage kills the leech")
	}
	if len(tb.G.Players[1].Hand) != 2 {
		t.Errorf("two kills: Cain loots 1 each time; hand %d", len(tb.G.Players[1].Hand))
	}
}
