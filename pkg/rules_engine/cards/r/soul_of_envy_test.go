package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSoulOfEnvy(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:    []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters:   []engine.CardRef{"chub", "leech"},
		BonusSouls: []engine.CardRef{"soul_of_envy"},
	})
	for range 2 {
		tb.G.Players[0].InPlay = append(tb.G.Players[0].InPlay, soulOf(tb, 0, "chub"))
	}
	tb.G.Players[2].InPlay = append(tb.G.Players[2].InPlay, soulOf(tb, 2, "chub"))
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "chub", foe) // Isaac's 3rd soul; Cain has the fewest
	settleAny(tb)
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("Cain's soul value = %d, want 1", s)
	}
}
