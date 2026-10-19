package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheEnigma(t *testing.T) {
	g, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: 2, Sets: sets, Characters: []engine.CardRef{"the_enigma", "cain"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, pl := range g.Players {
		if g.CardOf(pl.Character) == "amgine_eht" {
			t.Error("it never starts the game back side up")
		}
	}
	for _, ch := range g.CharactersOutside() {
		if ch == "amgine_eht" {
			t.Error("the back face is no character of its own")
		}
	}
}

func TestTheEnigmaCharacter(t *testing.T) { testCharacter(t, "the_enigma", "flip") }
