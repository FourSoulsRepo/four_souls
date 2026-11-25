package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Magic Marker is left out of every game until it works (19.96): no copy
// goes into a deck.
func TestMagicMarkerIsLeftOut(t *testing.T) {
	for seed := range uint64(20) {
		g, _, err := engine.NewGame(engine.Setup{Seed: seed, Players: 2, Sets: sets, BonusSouls: true})
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range g.Objects {
			if o.Card == "magic_marker" {
				t.Fatalf("seed %d: magic_marker is in the game, zone %+v", seed, o.Zone)
			}
		}
	}
}
