package b2

import "testing"

func TestMaskOfInfamy(t *testing.T) {
	tb := slayTable(t, "mask_of_infamy", seat("isaac"), seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	tb.G.Object(m).Damage = 3
	if dc := tb.G.Evasion(m); dc != 6 {
		t.Errorf("DC at 1 HP = %d, want 6", dc)
	}
}
