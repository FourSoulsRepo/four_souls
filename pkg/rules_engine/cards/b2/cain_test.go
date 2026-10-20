package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestCainGoesFirstWithoutARoll: the game's start events name Cain's
// player and the card, and no die is rolled (15.43).
func TestCainGoesFirstWithoutARoll(t *testing.T) {
	_, ev, err := engine.NewGame(engine.Setup{Seed: 7, Players: 3, Sets: []engine.CardSet{Set}, Characters: items("isaac", "cain", "eve")})
	if err != nil {
		t.Fatal(err)
	}
	var first []engine.Event
	for _, e := range ev {
		if e.Kind == engine.EvDiceRolled {
			t.Fatalf("a die was rolled: %+v", e)
		}
		if e.Kind == engine.EvFirstPlayer {
			first = append(first, e)
		}
	}
	if len(first) != 1 || first[0].Player != 1 || first[0].Card != "cain" || first[0].Text != "character" {
		t.Fatalf("first_player events %+v, want Cain's player 1", first)
	}
}

func TestCain(t *testing.T) {
	testExtraLoot(t, "cain")
	for _, chars := range [][]engine.CardRef{{"isaac", "cain"}, {"cain", "eve", "maggy"}} {
		g := newGame(t, chars...)
		want := engine.PlayerID(0)
		for i, c := range chars {
			if c == "cain" {
				want = engine.PlayerID(i)
			}
		}
		if g.Turn.Active != want {
			t.Errorf("%v: player %d goes first, want Cain's player %d", chars, g.Turn.Active, want)
		}
		// No roll: the game starts with Cain's player, and says why (15.43).
		if g.FirstRoll != nil {
			t.Errorf("%v: a first-player roll with Cain", chars)
		}
		if !hasItem(g, int(want), "sleight_of_hand") {
			t.Error("Cain has no Sleight of Hand")
		}
	}
}
