// Package world builds the games a bot may look at (ADR 033, Decisions
// 4.4 and 12.4 – 12.5): a copy of the real game, then for each world i a
// new game as one seat may imagine it (Game.Imagine).
//
// It is the one place that turns the real game into a bot's worlds. The
// arena (cmd/bot_arena, internal/arena) and the server's bot driver use
// it, so both give a bot the same thing. internal/bot never imports it:
// the bot gets only the World function this package returns, never the
// game it was made from.
package world

import (
	"fmt"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Func returns world i: a new imagined game each call. Equal i give
// equal games. It never changes the game it was made from.
type Func func(i int) (*engine.Game, error)

// For copies g now and returns the worlds of seat p: world i is the copy
// imagined by p with the seed mixed from seed and i, and the cards p
// remembers in known (a hard seat's Memory.Known; empty for others).
//
// The copy is made at once, so g may go on changing after For returns:
// the room makes the copy on its own goroutine and hands the worlds to
// a thinker (Decision 12.4). One Func is for one goroutine at a time.
func For(g *engine.Game, p engine.PlayerID, known engine.Known, seed uint64) (Func, error) {
	if p < 0 || int(p) >= len(g.Players) {
		return nil, fmt.Errorf("world: no seat %d", p)
	}
	snap, err := g.Clone()
	if err != nil {
		return nil, fmt.Errorf("world: copy: %w", err)
	}
	return func(i int) (*engine.Game, error) {
		return snap.Imagine(p, Seed(seed, i), known)
	}, nil
}

// Seed mixes a base seed and a world number into the seed of that world
// (SplitMix64), so neighbouring worlds share no pattern.
func Seed(seed uint64, i int) uint64 {
	z := seed + uint64(i+1)*0x9e3779b97f4a7c15 //nolint:gosec // i is a small world number, never negative
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
