package bot

import (
	"context"
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/world"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// newGame is a seeded game with every set.
func newGame(t *testing.T, seed uint64, players int) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: seed, Players: players, Sets: cards.Sets(), BonusSouls: true})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// input is seat p's input on g, with its own random source.
func input(t *testing.T, g *engine.Game, p engine.PlayerID, level Level, seed uint64) Input {
	t.Helper()
	w, err := world.For(g, p, engine.Known{}, seed)
	if err != nil {
		t.Fatal(err)
	}
	return Input{Seat: p, Level: level, Rand: rand.New(rand.NewPCG(seed, 1)), World: World(w)} //nolint:gosec // test moves
}

// TestChooseAllowed plays a few hundred steps with each level: every
// answer is accepted by the game, and the game's checksum does not move
// while the bot thinks.
func TestChooseAllowed(t *testing.T) {
	for _, level := range Levels() {
		g := newGame(t, 7, 3)
		for step := range 400 {
			if g.Over {
				break
			}
			p := g.Prompt().Who()[0]
			before, err := g.Checksum()
			if err != nil {
				t.Fatal(err)
			}
			in, err := Choose(context.Background(), input(t, g, p, level, uint64(step))) //nolint:gosec // step is small
			if err != nil {
				t.Fatalf("%s step %d: %v", level, step, err)
			}
			if after, err := g.Checksum(); err != nil || after != before {
				t.Fatalf("%s step %d: the bot changed the game", level, step)
			}
			if _, err := g.Apply(in); err != nil {
				t.Fatalf("%s step %d: %+v refused: %v", level, step, in, err)
			}
		}
	}
}

// TestChooseSameSeed: the same input gives the same move.
func TestChooseSameSeed(t *testing.T) {
	g := newGame(t, 3, 2)
	p := g.Prompt().Who()[0]
	for _, level := range Levels() {
		a, errA := Choose(context.Background(), input(t, g, p, level, 5))
		b, errB := Choose(context.Background(), input(t, g, p, level, 5))
		if errA != nil || errB != nil {
			t.Fatal(errA, errB)
		}
		if a.Kind != b.Kind || a.Choice != b.Choice || !slices.Equal(a.Objects, b.Objects) {
			t.Errorf("%s: %+v, then %+v", level, a, b)
		}
	}
}

// TestChooseErrors: no world, a failing world, a seat with nothing to
// do and an unknown level are errors, not moves.
func TestChooseErrors(t *testing.T) {
	g := newGame(t, 1, 2)
	p := g.Prompt().Who()[0]
	ctx := context.Background()
	if _, err := Choose(ctx, Input{Seat: p, Level: Random}); err == nil {
		t.Error("no world: no error")
	}
	bad := Input{Seat: p, Level: Random, Rand: rand.New(rand.NewPCG(1, 1)), World: func(int) (*engine.Game, error) { //nolint:gosec // test
		return nil, errors.New("broken")
	}}
	if _, err := Choose(ctx, bad); err == nil {
		t.Error("a failing world: no error")
	}
	idle := (p + 1) % engine.PlayerID(len(g.Players))
	if len(g.Prompt().Who()) == 1 {
		if _, err := Choose(ctx, input(t, g, idle, Random, 1)); err == nil {
			t.Error("a seat with nothing to do: no error")
		}
	}
	in := input(t, g, p, "clever", 1)
	if len(g.Allowed(p)) > 1 {
		if _, err := Choose(ctx, in); err == nil {
			t.Error("an unknown level: no error")
		}
	}
	if _, err := ParseLevel("clever"); err == nil {
		t.Error("ParseLevel took an unknown level")
	}
	for _, l := range Levels() {
		if got, err := ParseLevel(string(l)); err != nil || got != l {
			t.Errorf("ParseLevel(%q) = %q, %v", l, got, err)
		}
	}
}

// TestDiscards: up to maxDiscards different sets of n cards of the hand.
func TestDiscards(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2)) //nolint:gosec // test
	hand := []engine.ObjectID{4, 8, 15, 16, 23, 42}
	tmpl := engine.Intent{Player: 1, Kind: engine.IntentDiscard, Objects: hand}
	for _, n := range []int{1, 2, 5, 6} {
		got := discards(tmpl, n, rng)
		if len(got) == 0 || len(got) > maxDiscards {
			t.Fatalf("n %d: %d sets", n, len(got))
		}
		for i, in := range got {
			if len(in.Objects) != n || in.Kind != engine.IntentDiscard || in.Player != 1 {
				t.Errorf("n %d: bad set %+v", n, in)
			}
			for _, id := range in.Objects {
				if !slices.Contains(hand, id) {
					t.Errorf("n %d: %d is not in the hand", n, id)
				}
			}
			for _, other := range got[:i] {
				if slices.Equal(other.Objects, in.Objects) {
					t.Errorf("n %d: %v twice", n, in.Objects)
				}
			}
		}
	}
	if got := discards(tmpl, 6, rng); len(got) != 1 {
		t.Errorf("the whole hand: %d sets, want 1", len(got))
	}
	if got := discards(tmpl, 2, rng); len(got) != maxDiscards {
		t.Errorf("2 of 6: %d sets, want %d", len(got), maxDiscards)
	}
}

// TestOrders: the order now first, then up to maxOrders other orders of
// the same cards.
func TestOrders(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 3)) //nolint:gosec // test
	ids := []engine.ObjectID{3, 1, 2, 9}
	tmpl := engine.Intent{Player: 0, Kind: engine.IntentOrder, Objects: ids}
	got := orders(tmpl, rng)
	if len(got) != maxOrders+1 || !slices.Equal(got[0].Objects, ids) {
		t.Fatalf("orders: %v", got)
	}
	for i, in := range got {
		s := slices.Sorted(slices.Values(in.Objects))
		if !slices.Equal(s, []engine.ObjectID{1, 2, 3, 9}) {
			t.Errorf("order %v is not of the cards", in.Objects)
		}
		for _, other := range got[:i] {
			if slices.Equal(other.Objects, in.Objects) {
				t.Errorf("order %v twice", in.Objects)
			}
		}
	}
	if got := orders(engine.Intent{Objects: []engine.ObjectID{5, 6}}, rng); len(got) != 2 {
		t.Errorf("two cards: %d orders, want 2", len(got))
	}
}

// TestWeighted: a zero weight is never drawn; no positive weight draws
// any index.
func TestWeighted(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 4)) //nolint:gosec // test
	for range 200 {
		if i := weighted([]float64{0, 1, 0, 3}, rng); i != 1 && i != 3 {
			t.Fatalf("drew %d", i)
		}
		if i := weighted([]float64{0, 0}, rng); i < 0 || i > 1 {
			t.Fatalf("drew %d", i)
		}
	}
}

// TestWinBefore checks the fight odds against hand-worked cases.
func TestWinBefore(t *testing.T) {
	for _, c := range []struct {
		hits, misses int
		h, want      float64
	}{
		{1, 1, 0.5, 0.5},
		{2, 1, 0.5, 0.25},
		{1, 2, 0.5, 0.75},
		{3, 2, 0.5, 0.3125},
		{2, 2, 1, 1},
		{1, 1 << 10, 0.5, 1},
	} {
		if got := winBefore(c.hits, c.misses, c.h); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("winBefore(%d, %d, %v) = %v, want %v", c.hits, c.misses, c.h, got, c.want)
		}
	}
}
