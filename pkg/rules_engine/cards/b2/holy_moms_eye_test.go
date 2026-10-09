package b2

import "testing"

func TestHolyMomsEye(t *testing.T) {
	tb := slayTable(t, "holy_moms_eye", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "mystery_sack")
	if !tb.G.Object(tb.Find(0, "mystery_sack")).Charged {
		t.Error("the sack was not recharged")
	}
}
