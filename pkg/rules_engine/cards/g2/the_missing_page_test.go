package g2

import "testing"

func TestTheMissingPage(t *testing.T) {
	tb := table(t, seat("isaac", "the_butter_bean"), seat("cain", "the_missing_page"))
	tb.G.ForceRolls(5, 6)
	tb.Activate(0, "the_butter_bean", 0, "yes") // Cain forces the 5 to be rerolled: 6
	if !tb.G.Object(tb.Find(0, "the_butter_bean")).Charged {
		t.Error("the 5 should have been rerolled into a 6")
	}
	tb2 := table(t, seat("isaac", "the_butter_bean"), seat("cain", "the_missing_page"))
	tb2.G.ForceRolls(5, 6)
	tb2.Activate(0, "the_butter_bean", 0, "no")
	if tb2.G.Object(tb2.Find(0, "the_butter_bean")).Charged {
		t.Error("declined: the 5 stays")
	}
}
