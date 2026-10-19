package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestRedVise(t *testing.T) {
	tb := lootRoomTable(t, "red_vise", "a_penny")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentEndTurn})
	ev := tb.Do(engine.Intent{Player: 1, Kind: engine.IntentPass}) // the end phase starts
	if !enginetest.Has(ev, engine.EvDamaged, 0) {
		t.Error("Isaac takes 1 damage with 1 loot card in hand")
	}
}
