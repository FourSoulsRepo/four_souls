package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBackstabber(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "backstabber"), {Character: "cain", Souls: []engine.CardRef{"monstro"}}, seat("samson"),
	}})
	tb.G.ForceRolls(3, 3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe)
	if !tb.G.Players[1].Dead || len(tb.G.Players[0].Hand) != 1 {
		t.Error("Cain (more souls) can be attacked at 3+; killing him loots 1")
	}
	for _, o := range tb.Options() {
		if o == "player 3 (samson)" {
			t.Error("Samson has no more souls: not a target")
		}
	}
}
