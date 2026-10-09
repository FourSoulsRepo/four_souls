package b2

import (
	"strconv"
	"testing"
)

func TestBookOfBelial(t *testing.T) {
	cases := []struct {
		roll int
		mode string
		want int
	}{
		{3, "Add 1 to a roll.", 4},
		{3, "Subtract 1 from a roll.", 2},
		{6, "Add 1 to a roll.", 6},        // never above 6 (R-DICE-01)
		{1, "Subtract 1 from a roll.", 1}, // never below 1
	}
	for _, c := range cases {
		tb := rollTable(t, c.roll, seat("judas", "book_of_belial"), seat("cain"))
		tb.Pass(1)
		ev := tb.Activate(0, "book_of_belial", 0, c.mode, "roll of "+strconv.Itoa(c.roll))
		if got := resolvedRoll(t, ev); got != c.want {
			t.Errorf("%d, %s: resolved as %d, want %d", c.roll, c.mode, got, c.want)
		}
	}
	tb := rollTable(t, 3, seat("judas", "book_of_belial"), seat("cain"))
	tb.Pass(1)
	tb.Activate(0, "book_of_belial", 0, "cancel")
	if !tb.G.Object(tb.Find(0, "book_of_belial")).Charged {
		t.Error("cancelling the mode choice paid the cost")
	}
}
