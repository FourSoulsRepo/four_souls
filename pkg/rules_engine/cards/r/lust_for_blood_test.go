package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLustForBlood(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain", "spoon_bender")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "lust_for_blood")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe) // Lust for Blood! fills the slot
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err == nil {
		t.Fatal("Isaac must attack Cain first")
	}
	tb.G.ForceRolls(4, 4)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe, "spoon_bender")
	if !tb.G.Players[1].Dead || !hasItem(tb.G, 0, "spoon_bender") {
		t.Error("Cain dies at 4+ and gives Isaac the item he would destroy")
	}
}
