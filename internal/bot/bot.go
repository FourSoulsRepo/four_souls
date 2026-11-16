// Package bot picks the moves of a bot seat (ADR 030, Decision 4).
//
// A bot never holds the real game. Choose gets an Input whose World
// returns games as the seat may imagine them (Game.Imagine, built by
// internal/world); every game the bot reads or copies comes from there.
// The package imports the engine and the standard library only, and its
// code names no card: TestBotRules checks both.
package bot

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Level is how a bot plays (Decision 1).
type Level string

// The levels so far. Random and Policy are for tests and the arena only
// (Decision 1.2); the server refuses them.
const (
	// Random picks one of the candidates at random (Decision 5).
	Random Level = "random"
	// Policy plays by the cheap playout policy alone (Decision 6).
	Policy Level = "policy"
	// Standard tries each candidate on one imagined world and plays the
	// best score, with one mistake in ten (Decision 8).
	Standard Level = "standard"
	// Hard plays each candidate out on many imagined worlds to the end
	// of the turn, halving the candidates round by round, with the
	// opponents weighted by strength; no mistakes (Decision 9).
	Hard Level = "hard"
	// Flat is standard with every loot card and item worth alike, as
	// before the measured card values (Decision 20): the arena measures
	// the table against it, and make bot-values plays its positions
	// with it. For the arena and tools only, as Random and Policy.
	Flat Level = "flat"
)

// Levels lists every level Choose knows, in a fixed order.
func Levels() []Level { return []Level{Random, Policy, Standard, Hard, Flat} }

// ParseLevel returns the level named s.
func ParseLevel(s string) (Level, error) {
	for _, l := range Levels() {
		if string(l) == s {
			return l, nil
		}
	}
	return "", fmt.Errorf("bot: unknown level %q", s)
}

// World returns world i: a new game as the seat may imagine it. Equal i
// give equal games. internal/world builds it from the real game.
type World func(i int) (*engine.Game, error)

// Budget bounds the hard level's search: whichever ends first
// (Decision 9.8). Playouts 0 sets no cap on playouts, Time 0 no time
// limit; with both 0 it thinks for 1.5 s. A playout count alone
// repeats: the arena uses it (Decision 9.10). The other levels ignore
// the budget.
type Budget struct {
	Playouts int
	Time     time.Duration
}

// Input is everything a bot gets for one move (Decision 4.3). It holds
// no game: only World, which imagines one.
type Input struct {
	Seat  engine.PlayerID
	Level Level
	// Rand is the bot's own random source, never the game's RNG
	// (Decision 4.7).
	Rand   *rand.Rand
	Budget Budget
	World  World
}

// Choose returns the seat's move: one of Allowed, completed for a
// discard or an order (Decision 4.5). With one allowed intent it
// answers at once (Decision 4.6), unless it is a discard of some of the
// hand: then the levels that score weigh the sets as any move, so a
// forced discard drops a cheap card (Decision 20). An error means no
// world or nothing allowed; the caller then uses its forced answer
// (Decision 12.10).
// ctx ends the hard level's search early; it then plays the best so far
// (Decision 9.9).
func Choose(ctx context.Context, in Input) (engine.Intent, error) {
	if in.World == nil || in.Rand == nil {
		return engine.Intent{}, errors.New("bot: input has no world or no random source")
	}
	w, err := in.World(0)
	if err != nil {
		return engine.Intent{}, fmt.Errorf("bot: world: %w", err)
	}
	allowed := w.Allowed(in.Seat)
	if len(allowed) == 0 {
		return engine.Intent{}, fmt.Errorf("bot: seat %d has nothing to do", in.Seat)
	}
	if len(allowed) == 1 && (!partDiscard(w, in.Seat, allowed[0]) || in.Level == Random || in.Level == Policy) {
		return complete(w, in.Seat, allowed[0], in.Rand), nil
	}
	switch in.Level {
	case Random:
		cs := candidates(w, in.Seat, allowed, in.Rand)
		return cs[in.Rand.IntN(len(cs))], nil
	case Policy:
		return policyPick(w, in.Seat, allowed, in.Rand), nil
	case Standard, Flat:
		return standardPick(w, in.Seat, allowed, in.Rand, mistakeOdds, in.Level == Flat)
	case Hard:
		return hardPick(ctx, in, w, allowed)
	default:
		return engine.Intent{}, fmt.Errorf("bot: unknown level %q", in.Level)
	}
}

// partDiscard: in is a discard of seat p on world w that keeps some of
// the cards, so which ones is a choice.
func partDiscard(w *engine.Game, p engine.PlayerID, in engine.Intent) bool {
	return in.Kind == engine.IntentDiscard && w.PromptFor(p).Count < len(in.Objects)
}
