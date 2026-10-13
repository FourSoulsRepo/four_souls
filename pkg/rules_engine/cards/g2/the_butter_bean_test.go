package g2

import "testing"

func TestTheButterBean(t *testing.T) {
	tb := table(t, seat("isaac", "the_butter_bean"), seat("cain"))
	tb.G.ForceRolls(6)
	tb.Activate(0, "the_butter_bean", 0)
	if !tb.G.Object(tb.Find(0, "the_butter_bean")).Charged {
		t.Error("a 6 recharges it")
	}
	tb.G.ForceRolls(3)
	tb.Activate(0, "the_butter_bean", 0)
	if tb.G.Object(tb.Find(0, "the_butter_bean")).Charged {
		t.Error("a 3 does not recharge it")
	}
}
