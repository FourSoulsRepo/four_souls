package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The rooms option at game start (R-SETUP-16); the room rules themselves
// (R-TURN-12, R-TURN-17, R-CARD-33) are situation files.

func TestRoomsOption(t *testing.T) {
	for _, on := range []bool{false, true} {
		g, _, err := engine.NewGame(engine.Setup{Seed: 3, Players: 2, Sets: sets, Rooms: on})
		if err != nil {
			t.Fatal(err)
		}
		if got := len(g.Rooms) > 0; got != on {
			t.Errorf("rooms %v: a room slot = %v", on, got)
		}
		for _, id := range g.Decks[engine.LootDeck] {
			if g.Def(id).Kind == engine.RoomCard {
				t.Fatal("a room card in the loot deck")
			}
		}
		if !on && len(g.Decks[engine.RoomDeck]) > 0 {
			t.Error("without the option, rooms stay out of the game (R-SETUP-16)")
		}
		// Every viewer sees the room slot (14.13).
		v := g.View(engine.Viewer{Kind: engine.ViewSpectator})
		if got := len(v.Rooms) == 1 && v.Rooms[0].Top != nil; got != on {
			t.Errorf("rooms %v: the view shows a room = %v", on, got)
		}
	}
}
