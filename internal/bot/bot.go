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
)

// Levels lists every level Choose knows, in a fixed order.
func Levels() []Level { return []Level{Random, Policy, Standard} }

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

// Budget bounds a search: whichever ends first (Decision 9.8). The
// levels so far do not search and ignore it.
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
// answers at once (Decision 4.6). An error means no world or nothing
// allowed; the caller then uses its forced answer (Decision 12.10).
func Choose(_ context.Context, in Input) (engine.Intent, error) {
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
	if len(allowed) == 1 {
		return complete(w, in.Seat, allowed[0], in.Rand), nil
	}
	switch in.Level {
	case Random:
		cs := candidates(w, in.Seat, allowed, in.Rand)
		return cs[in.Rand.IntN(len(cs))], nil
	case Policy:
		return policyPick(w, in.Seat, allowed, in.Rand), nil
	case Standard:
		return standardPick(w, in.Seat, allowed, in.Rand, mistakeOdds)
	default:
		return engine.Intent{}, fmt.Errorf("bot: unknown level %q", in.Level)
	}
}
