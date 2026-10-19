package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBloodDonation(t *testing.T) {
	tb := roomTable(t, "blood_donation")
	useRoom(tb, 0, 0)
	if tb.G.Players[0].Cents != 3 || tb.G.Players[0].Damage != 1 {
		t.Errorf("cents %d, damage %d; want 3 and 1", tb.G.Players[0].Cents, tb.G.Players[0].Damage)
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 1, Kind: engine.IntentActivate, Objects: []engine.ObjectID{room(tb)}}); err == nil {
		t.Error("only the active player uses a room (R-CARD-36)")
	}
}
