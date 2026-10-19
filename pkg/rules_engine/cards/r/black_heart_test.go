package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBlackHeartPreventsOne(t *testing.T) {
	tb := lootTable(t, "black_heart", "mega_bomb")
	tb.G.Turn.LootPlays = 2
	tb.Play(0, "black_heart", foe)
	tb.Play(0, "mega_bomb", foe)
	if got := tb.G.Players[1].Damage; got != 1 {
		t.Errorf("Cain's damage = %d, want 2 - 1 prevented", got)
	}
}

func TestBlackHeartWhenThatPlayerDies(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"black_heart", "mega_bomb", "mega_bomb"}}, seat("cain"), seat("samson"),
	}})
	tb.G.Turn.LootPlays = 3
	tb.Play(0, "black_heart", foe)
	tb.Play(0, "mega_bomb", foe) // 1 prevented, 1 taken
	tb.Play(0, "mega_bomb", foe) // Cain dies
	if !tb.G.Players[1].Dead {
		t.Fatal("Cain should be dead")
	}
	if got := tb.G.Players[2].Damage; got != 2 {
		t.Errorf("Samson's damage = %d, want 2", got)
	}
	if got := tb.G.Players[0].Damage; got != 0 {
		t.Errorf("the player of Black Heart took %d damage, want none", got)
	}
}
