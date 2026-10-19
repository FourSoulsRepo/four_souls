package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Tests of the room rules (R-SETUP-16, R-TURN-12, R-TURN-17, R-CARD-36).

func TestRoomStep(t *testing.T) {
	tb := roomTable(t, "blessing_of_greed")
	tb.G.PutOnDeck(engine.RoomDeck, "dice_room_5")
	tb.Attack("fly", 6) // a monster died this turn
	tb.EndTurn("blessing_of_greed")
	if !inDiscard(tb.G, engine.RoomDeck, "blessing_of_greed") || tb.G.Object(room(tb)).Card != "dice_room_5" {
		t.Fatal("the active player put the room into discard; the slot was refilled")
	}
	if !inDiscard(tb.G, engine.MonsterDeck, "leech") {
		t.Error("the new room's \"enters play\" ability happens in the end phase (R-TURN-18)")
	}
	if tb.G.Players[1].Cents != 0 {
		t.Error("Blessing of Greed left before Cain's turn")
	}
}

func TestRoomStepNeedsADeath(t *testing.T) {
	tb := roomTable(t, "blessing_of_greed")
	tb.EndTurn() // no question: no monster died
	if tb.G.Object(room(tb)).Card != "blessing_of_greed" {
		t.Error("the room stays")
	}
}

// A Dice Room putting itself into discard is not the room step: the
// active player may still replace the next room (R-TURN-17).
func TestDiceRoomIsNotTheRoomStep(t *testing.T) {
	tb := roomTable(t, "dice_room_3")
	tb.G.PutOnDeck(engine.RoomDeck, "blessing_of_greed")
	tb.G.PutOnDeck(engine.RoomDeck, "spider_webs")
	tb.Attack("fly", 6)
	tb.EndTurn("spider_webs")
	if got := tb.G.Object(room(tb)).Card; got != "blessing_of_greed" {
		t.Errorf("room = %s, want blessing_of_greed", got)
	}
}

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
	}
}

func TestOnlyTheActivePlayerUsesRooms(t *testing.T) {
	tb := roomTable(t, "blood_donation")
	if _, err := tb.G.Apply(engine.Intent{Player: 1, Kind: engine.IntentActivate, Objects: []engine.ObjectID{room(tb)}}); err == nil {
		t.Error("Cain used the room on Isaac's turn (R-CARD-36)")
	}
}
