package r

import "testing"

func TestPlanetarium(t *testing.T) {
	tb := roomTable(t, "planetarium")
	tb.EndTurn("ask.yes")
	if len(itemsOf(tb.G, 1)) != 1 {
		t.Error("Cain may gain +1 treasure at the start of his turn")
	}
}
