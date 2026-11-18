package r

import "testing"

func TestBloodDonation(t *testing.T) {
	tb := roomTable(t, "blood_donation")
	useRoom(tb, 0, 0)
	if tb.G.Players[0].Cents != 3 || tb.G.Players[0].Damage != 1 {
		t.Errorf("cents %d, damage %d; want 3 and 1", tb.G.Players[0].Cents, tb.G.Players[0].Damage)
	}
}
