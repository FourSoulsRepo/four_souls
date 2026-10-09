package b2

import "testing"

func TestPortableSlotMachine(t *testing.T) {
	tb := itemTable(t, items("portable_slot_machine"))
	tb.G.Players[0].Cents = 3
	tb.G.ForceRolls(1)
	tb.Activate(0, "portable_slot_machine", 0)
	if c, h := tb.G.Players[0].Cents, len(tb.G.Players[0].Hand); c != 0 || h != 1 {
		t.Errorf("cents %d hand %d, want 0 and 1", c, h)
	}
}
