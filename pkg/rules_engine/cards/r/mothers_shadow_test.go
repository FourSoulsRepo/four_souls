package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMothersShadow(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain", "spoon_bender"), seat("samson", "jawbone")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "mothers_shadow")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "spoon_bender", "jawbone")
	if !tb.G.Players[1].Dead || !tb.G.Players[2].Dead || tb.G.Players[0].Dead {
		t.Error("Cain and Samson tie for the most items: both die")
	}
}
