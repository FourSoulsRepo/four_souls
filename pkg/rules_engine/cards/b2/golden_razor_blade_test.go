package b2

import "testing"

func TestGoldenRazorBlade(t *testing.T) {
	tb := itemTable(t, items("golden_razor_blade"))
	tb.G.Players[0].Cents = 10
	tb.Activate(0, "golden_razor_blade", 0, foe)
	tb.Activate(0, "golden_razor_blade", 0, foe) // paid: usable while deactivated
	if !tb.G.Players[1].Dead || tb.G.Players[0].Cents != 0 {
		t.Errorf("dead %v cents %d", tb.G.Players[1].Dead, tb.G.Players[0].Cents)
	}
}
