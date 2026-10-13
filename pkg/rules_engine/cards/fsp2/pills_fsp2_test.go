package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPillsFsp2Cancels(t *testing.T) {
	plays := 2
	tb := setup(t, engine.SituationSetup{
		Players:   []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"bomb", "pills_fsp2"}}, seat("cain")},
		Monsters:  []engine.CardRef{"fly", "leech"},
		LootPlays: &plays,
	})
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "bomb")}}); err != nil {
		t.Fatal(err)
	}
	tb.Choose("fly")
	tb.G.ForceRolls(1)
	tb.Play(0, "pills_fsp2", "bomb")
	if top(tb, 0) != "fly" {
		t.Error("the bomb should be cancelled")
	}
}

func TestPillsFsp2Discards(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"pills_fsp2"}}, {Character: "cain", Hand: []engine.CardRef{"bomb"}},
	}})
	tb.G.ForceRolls(3)
	tb.Play(0, "pills_fsp2", "bomb")
	if n := len(tb.G.Players[1].Hand); n != 0 {
		t.Errorf("Cain's hand = %d, want 0", n)
	}
}
