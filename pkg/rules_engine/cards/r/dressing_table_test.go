package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDressingTable(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "spoon_bender", "jawbone"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "dressing_table")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "spoon_bender", "done")
	if hasItem(tb.G, 0, "spoon_bender") || !hasItem(tb.G, 0, "jawbone") || len(tb.G.Players[0].InPlay) != 2 {
		t.Error("the Spoon Bender is rerolled; the Jawbone stays")
	}
}
