package r

import "testing"

func TestAllHallowsEve(t *testing.T) {
	tb := roomTable(t, "all_hallows_eve")
	if tb.G.CanGainSouls(0) || tb.G.CanGainSouls(1) {
		t.Error("players can't gain souls")
	}
}
