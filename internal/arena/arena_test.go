package arena

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/bot"
)

// TestNoRefusedIntent plays 30 whole policy games, every set, 2 to 4
// players: the engine refuses no bot intent (ADR 030, Decision 18.8).
// The games run side by side, ten per player count.
func TestNoRefusedIntent(t *testing.T) {
	for players := 2; players <= 4; players++ {
		levels := make([]bot.Level, players)
		for i := range levels {
			levels[i] = bot.Policy
		}
		for k := range 10 {
			m := Match{Levels: levels, Games: 1, Seed: uint64(100*players + k)} //nolint:gosec // small numbers
			t.Run(fmt.Sprintf("%dp-seed%d", players, m.Seed), func(t *testing.T) {
				t.Parallel()
				r, err := Play(context.Background(), m)
				if err != nil {
					t.Fatal(err)
				}
				for _, x := range r.Refused {
					t.Error("refused:", x)
				}
				if r.Games != 1 || r.Moves == 0 {
					t.Errorf("played %d games, %d moves", r.Games, r.Moves)
				}
			})
		}
	}
}

// TestRandomLevel: random seats play whole games too, and the result
// counts every seat and game.
func TestRandomLevel(t *testing.T) {
	m := Match{Levels: []bot.Level{bot.Random, bot.Policy}, Games: 2, Seed: 1, Cap: 4}
	r, err := Play(context.Background(), m)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Refused) > 0 {
		t.Error(r.Refused)
	}
	if r.Games != 2 || r.Seats[bot.Random] != 2 || r.Seats[bot.Policy] != 2 {
		t.Errorf("games %d, seats %v", r.Games, r.Seats)
	}
	if r.Finished+len(r.Stalls) != r.Games {
		t.Errorf("%d finished and %d stalls of %d games", r.Finished, len(r.Stalls), r.Games)
	}
	if len(r.Turns) != 2 || r.MaxTurns() > 4 {
		t.Errorf("turns %v past the cap of 4", r.Turns)
	}
	if r.Passes == 0 || r.Passes > r.Steps {
		t.Errorf("%d passes of %d intents", r.Passes, r.Steps)
	}
	if m.Name() != "random-policy" {
		t.Errorf("name %q", m.Name())
	}
}

// TestPlayBadMatch: a line-up of one seat is an error.
func TestPlayBadMatch(t *testing.T) {
	if _, err := Play(context.Background(), Match{Levels: []bot.Level{bot.Random}, Games: 1}); err == nil {
		t.Error("one seat: no error")
	}
}

// TestGuards: a game turn or a game over its intents is ended and listed
// with its seed.
func TestGuards(t *testing.T) {
	for _, m := range []Match{
		{Levels: []bot.Level{bot.Policy, bot.Policy}, Games: 1, Seed: 3, TurnIntents: 5},
		{Levels: []bot.Level{bot.Policy, bot.Policy}, Games: 1, Seed: 3, GameIntents: 50},
	} {
		r, err := Play(context.Background(), m)
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Guards) != 1 || !strings.HasPrefix(r.Guards[0], "game 3:") {
			t.Errorf("guards %v", r.Guards)
		}
		if r.Finished != 0 || r.Steps > 50 {
			t.Errorf("finished %d, %d intents", r.Finished, r.Steps)
		}
	}
}
