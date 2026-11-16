package world

import (
	"bytes"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

func newGame(t *testing.T) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 11, Players: 3, Sets: cards.Sets(), BonusSouls: true})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func save(t *testing.T, g *engine.Game) []byte {
	t.Helper()
	b, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// world is world i, or the test fails.
func world(t *testing.T, worlds Func, i int) *engine.Game {
	t.Helper()
	w, err := worlds(i)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

// TestFor: world i is Imagine of the game for the seat, the same for the
// same i, other for another i; the game never changes, and a change to
// the game after For does not reach the worlds.
func TestFor(t *testing.T) {
	g := newGame(t)
	before := save(t, g)
	const p = engine.PlayerID(1)
	worlds, err := For(g, p, engine.Known{}, 42)
	if err != nil {
		t.Fatal(err)
	}
	want, err := g.Imagine(p, Seed(42, 0), engine.Known{})
	if err != nil {
		t.Fatal(err)
	}
	w0, err := worlds(0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(save(t, w0), save(t, want)) {
		t.Error("world 0 is not the game imagined with its seed")
	}
	again := world(t, worlds, 0)
	if !bytes.Equal(save(t, again), save(t, w0)) {
		t.Error("world 0 twice differs")
	}
	w1 := world(t, worlds, 1)
	if bytes.Equal(save(t, w1), save(t, w0)) {
		t.Error("worlds 0 and 1 are the same")
	}
	if !bytes.Equal(save(t, g), before) {
		t.Fatal("For changed the game")
	}
	// Play on the game: the worlds stay those of the moment of For.
	q := g.Prompt().Who()[0]
	in := g.Allowed(q)[0]
	if in.Kind == engine.IntentDiscard {
		in.Objects = in.Objects[:g.PromptFor(q).Count]
	}
	if _, err := g.Apply(in); err != nil {
		t.Fatal(err)
	}
	later := world(t, worlds, 0)
	if !bytes.Equal(save(t, later), save(t, w0)) {
		t.Error("a move on the game changed world 0")
	}
}

// TestForBadSeat: a seat the game does not have is an error.
func TestForBadSeat(t *testing.T) {
	g := newGame(t)
	for _, p := range []engine.PlayerID{-1, 3} {
		if _, err := For(g, p, engine.Known{}, 1); err == nil {
			t.Errorf("seat %d: no error", p)
		}
	}
}

// TestSeed: different worlds and base seeds give different seeds.
func TestSeed(t *testing.T) {
	seen := map[uint64]bool{}
	for base := range uint64(20) {
		for i := range 50 {
			s := Seed(base, i)
			if seen[s] {
				t.Fatalf("seed %d world %d repeats", base, i)
			}
			seen[s] = true
		}
	}
}
