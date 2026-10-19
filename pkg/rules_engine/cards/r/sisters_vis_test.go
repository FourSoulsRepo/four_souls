package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSistersVis(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"sisters_vis"},
	})
	tb.G.Object(tb.G.Monsters[0].Cards[0]).Damage = 4 // 2 HP left
	tb.G.ForceRolls(6, 6)                             // Isaac hits, then Samson hits and kills
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "sisters_vis", "player 3 (samson)")
	if tb.G.SoulValue(0) != 1 {
		t.Error("the active player gains the soul")
	}
	if len(tb.G.Players[2].InPlay) != 1 {
		t.Error("Samson killed it with Isaac: he gains +1 treasure")
	}
}
