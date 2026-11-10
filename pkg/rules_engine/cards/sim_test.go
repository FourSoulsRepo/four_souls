package cards

import (
	"math/rand/v2"
	"os"
	"strconv"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
)

// maxSimSteps caps one simulated game.
const maxSimSteps = 6000

// TestSimulations plays whole seeded games with random legal moves for
// 2, 3 and 4 players and checks the invariants after every step (5.9).
// Replays are checked by FuzzSets. CI plays a few; locally:
//
//	FOUR_SOULS_SIMS=3000 go test -run TestSimulations ./cards/
func TestSimulations(t *testing.T) {
	games := 100
	if n, err := strconv.Atoi(os.Getenv("FOUR_SOULS_SIMS")); err == nil && n > 0 {
		games = n
	}
	finished, steps := 0, 0
	for i := range games {
		seed := uint64(i + 1)
		players := 2 + i%3
		s := engine.Setup{Seed: seed, Players: players, Sets: Sets(), BonusSouls: true, Rooms: seed%2 == 1, DeckRatios: seed%3 == 0}
		g, n := simulate(t, s)
		steps += n
		if g.Over {
			finished++
		}
	}
	t.Logf("%d games, %d finished with a winner, %d steps", games, finished, steps)
	// Krampus at every difficulty (14.22): a quarter as many games.
	won, lost := 0, 0
	for i := range max(games/4, 3) {
		g, n := simulate(t, krampusSetup(uint64(i+1), krampusLevels[i%3]))
		steps += n
		switch {
		case g.Over && g.Challenge.Lost:
			lost++
		case g.Over:
			won++
		}
	}
	t.Logf("Krampus: %d won, %d lost", won, lost)
}

// krampusLevels are the Krampus difficulties (R-KRAMPUS-02).
var krampusLevels = []string{engine.DifficultyNormal, engine.DifficultyHard, engine.DifficultyUltra}

// krampusSetup is a Krampus challenge game (ADR 028, 14.22).
func krampusSetup(seed uint64, difficulty string) engine.Setup {
	return engine.Setup{
		Seed: seed, Players: 2, Sets: append(Sets(), challenge.Set), Rooms: seed%2 == 1,
		Challenge: &challenge.Krampus, Difficulty: difficulty,
	}
}

// simulate plays one game; it returns the game and its length.
func simulate(t *testing.T, s engine.Setup) (*engine.Game, int) {
	t.Helper()
	seed := s.Seed
	g, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	rng := rand.New(rand.NewPCG(seed, 99)) //nolint:gosec // test moves, not security
	step := 0
	for ; step < maxSimSteps; step++ {
		if g.Over {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		if len(allowed) == 0 {
			t.Fatalf("seed %d step %d: nothing allowed, prompt %+v", seed, step, p)
		}
		in := pick(rng, allowed)
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatalf("seed %d step %d: allowed %+v refused: %v", seed, step, in, err)
		}
		if err := g.CheckInvariants(); err != nil {
			t.Fatalf("seed %d step %d after %+v: %v", seed, step, in, err)
		}
	}
	return g, step
}

// pick prefers doing something over passing, so games move on.
func pick(rng *rand.Rand, allowed []engine.Intent) engine.Intent {
	var acts []engine.Intent
	for _, in := range allowed {
		if in.Kind != engine.IntentPass && in.Kind != engine.IntentEndTurn {
			acts = append(acts, in)
		}
	}
	if len(acts) > 0 && rng.IntN(3) > 0 {
		return acts[rng.IntN(len(acts))]
	}
	return allowed[rng.IntN(len(allowed))]
}
