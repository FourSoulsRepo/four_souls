package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestDelirium(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain")}, Monsters: items("delirium", "fly")}, Set)
	fly, _ := tb.G.Monsters[1].TopOf()
	if dc := tb.G.Evasion(fly); dc != 3 {
		t.Errorf("fly DC = %d, want 3", dc)
	}
	killWithSixes(t, tb)
	// 6th from the top, then the empty slot drew the top card: now 5th.
	if top := tb.G.DeckTop(engine.MonsterDeck, 6); tb.G.Object(top[4]).Card != "delirium" {
		t.Error("Delirium is not 6th from the top")
	}
}
