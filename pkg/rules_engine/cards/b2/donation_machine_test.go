package b2

import "testing"

func TestDonationMachine(t *testing.T) {
	tb := itemTable(t, items("donation_machine", "breakfast", "the_d6"))
	tb.Start(0, "donation_machine", 0)
	if opts := tb.Options(); len(opts) != 2 || opts[0] != "breakfast" {
		t.Fatalf("options %v, want the breakfast only (not this, not eternal)", opts)
	}
	tb.Settle("breakfast", foe)
	if !hasItem(tb.G, 1, "breakfast") || tb.G.Players[0].Cents != 8 {
		t.Error("not given for 8¢")
	}
}
