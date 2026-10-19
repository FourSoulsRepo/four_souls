package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGus(t *testing.T) {
	tb := roomTable(t, "gus")
	tb.G.Object(room(tb)).Damage = 8
	tb.G.ForceRolls(6, 6) // Cain rolls 6 to help; Isaac's attack roll hits
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "gus")
	settleAny(tb)
	if tb.G.SoulValue(0) != 1 {
		t.Errorf("Isaac killed Gus and gains its soul: soul value %d", tb.G.SoulValue(0))
	}
	if !tb.G.Attack.On && len(tb.G.Rooms[0].Cards) == 0 {
		t.Error("the room slot is refilled")
	}
}
