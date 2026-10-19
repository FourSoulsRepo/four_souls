package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCurseOfEmpathy(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{}, Hand: []engine.CardRef{"bomb"}, Cents: 3}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	giveCurse(tb, 0, "curse_of_empathy")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "bomb")
	if handHas(tb.G, 0, "bomb") || tb.G.Players[0].Cents != 3 {
		t.Errorf("a monster died: discard and lose 1¢ (+1¢ from the fly); cents %d", tb.G.Players[0].Cents)
	}
}
