package arena

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/bot"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// TestNoRefusedIntent plays 30 whole policy games, every set, 2 to 4
// players: the engine refuses no bot intent (ADR 033, Decision 18.8).
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
// counts every seat and game, and each level's moves and turns.
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
	moves, turns := 0, 0
	for l, c := range r.Levels {
		moves += c.Moves
		turns += c.Turns
		if c.NoAttack > c.Turns || c.LootPlays > c.Moves || c.Purchases > c.Moves {
			t.Errorf("%s: %+v", l, c)
		}
	}
	if moves != r.Moves {
		t.Errorf("levels made %d moves of %d", moves, r.Moves)
	}
	// A game cut by the cap ends after turn 4; each turn counts once.
	if turns != 8 {
		t.Errorf("levels had %d turns, want 8 (turns %v)", turns, r.Turns)
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

// TestStandardNoLoop replays the games where standard seats once used
// the same items again and again in one turn, past the turn guard: a
// move that gains nothing over doing nothing is not made (ADR 033,
// Decision 8.6). Each game is played to its end or the cap.
//
// It runs only with FOUR_SOULS_ARENA_GAMES set, never with -short:
// make bots runs it after the arena. Its games are whole standard games
// as long as the bots make them: 3-player seed 28 lasts 129 turns and
// took 15 minutes alone (2026-11-17), past go test's 10-minute timeout.
// Plain go test checks the rule on built tables in internal/bot:
// TestJudge (a move must beat doing nothing by wAct and its noise),
// TestStandardKeepsUselessLoot and TestStandardNoHopelessAttack; and
// TestGuards here checks the loop guards themselves.
func TestStandardNoLoop(t *testing.T) {
	if testing.Short() || os.Getenv("FOUR_SOULS_ARENA_GAMES") == "" {
		t.Skip("whole standard games: FOUR_SOULS_ARENA_GAMES=1, or make bots")
	}
	s, r := bot.Standard, bot.Random
	for _, m := range []Match{
		{Levels: []bot.Level{r, s}, Seed: 10},
		{Levels: []bot.Level{r, s}, Seed: 38},
		{Levels: []bot.Level{s, s}, Seed: 10},
		{Levels: []bot.Level{s, s, s}, Seed: 28},
		{Levels: []bot.Level{s, s, s}, Seed: 41},
	} {
		m.Games = 1
		t.Run(fmt.Sprintf("%s-seed%d", m.Name(), m.Seed), func(t *testing.T) {
			t.Parallel()
			res, err := Play(context.Background(), m)
			if err != nil {
				t.Fatal(err)
			}
			for _, x := range res.Guards {
				t.Error("loop guard:", x)
			}
			for _, x := range res.Refused {
				t.Error("refused:", x)
			}
		})
	}
}

// TestMemories: each hard seat gets a memory of its own seat, the
// other levels none; NoMemory gives none (Decision 19.13).
func TestMemories(t *testing.T) {
	mems := memories([]bot.Level{bot.Standard, bot.Hard, bot.Policy, bot.Hard}, false)
	for i, want := range []bool{false, true, false, true} {
		if got := mems[i] != nil; got != want {
			t.Errorf("seat %d: memory %v, want %v", i, got, want)
		}
	}
	for i, m := range memories([]bot.Level{bot.Hard, bot.Hard}, true) {
		if m != nil {
			t.Errorf("seat %d has a memory with NoMemory", i)
		}
	}
}

// TestChallenges: each challenge at each difficulty, bot with bot
// (Decision 15.5), refuses no bot intent and trips no loop guard:
// random seats play whole games, standard seats the first turns.
func TestChallenges(t *testing.T) {
	for _, c := range cards.Challenges() {
		for _, l := range c.Levels {
			for _, m := range []Match{
				{Levels: []bot.Level{bot.Random, bot.Random}, Games: 2},
				{Levels: []bot.Level{bot.Standard, bot.Hard}, Games: 1, Cap: 2, Playouts: 8},
			} {
				m.Challenge, m.Difficulty, m.Seed = c, l.Difficulty, 3
				t.Run(m.Name(), func(t *testing.T) {
					t.Parallel()
					r, err := Play(context.Background(), m)
					if err != nil {
						t.Fatal(err)
					}
					for _, x := range append(r.Refused, r.Guards...) {
						t.Error(x)
					}
					if r.Finished > 0 && r.Wins[m.Levels[0]] == 0 && r.Lost == 0 {
						t.Errorf("%d ended, none won or lost", r.Finished)
					}
				})
			}
		}
	}
}
