package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/arena"
	"github.com/FourSoulsRepo/four_souls/internal/bot"
)

// TestRun plays two short games: the report names the line-up, the
// wins by level and the stalls; a game cut by the cap fails no gate in
// a line-up of random or policy seats.
func TestRun(t *testing.T) {
	for _, match := range []string{"random,policy", "policy,policy"} {
		var out bytes.Buffer
		if err := run(context.Background(), []string{"-games", "2", "-cap", "2", "-match", match}, &out); err != nil {
			t.Errorf("%s: err %v, want no gate", match, err)
		}
		for _, want := range []string{"2 games", "policy ", "stall: seed 1", "for reference"} {
			if !strings.Contains(out.String(), want) {
				t.Errorf("%s: no %q in:\n%s", match, want, out.String())
			}
		}
		if strings.Contains(out.String(), "FAIL") {
			t.Errorf("%s: a gate failed:\n%s", match, out.String())
		}
	}
}

// TestReportGates: a guard hit and a refused intent fail any line-up;
// the finish share fails only a finish-gated one.
func TestReportGates(t *testing.T) {
	base := arena.Result{
		Match: arena.Match{Levels: []bot.Level{bot.Policy, bot.Policy}, Cap: 150},
		Games: 2, Finished: 0, Wins: map[bot.Level]int{}, Seats: map[bot.Level]int{bot.Policy: 4},
	}
	if !report(&printer{w: &bytes.Buffer{}}, base, 0) {
		t.Error("policy line-up with no winner failed a gate")
	}
	r := base
	r.Guards = []string{"game 1: turn 3 over 1000 intents"}
	var out bytes.Buffer
	if report(&printer{w: &out}, r, 0) || !strings.Contains(out.String(), "loop guard: game 1") {
		t.Errorf("a guard hit passed:\n%s", out.String())
	}
	r = base
	r.Refused = []string{"game 1 step 2"}
	if report(&printer{w: &bytes.Buffer{}}, r, 0) {
		t.Error("a refused intent passed")
	}
	finishGated[bot.Policy] = true
	defer delete(finishGated, bot.Policy)
	if report(&printer{w: &bytes.Buffer{}}, base, 0) {
		t.Error("a gated line-up with no winner passed")
	}
}

// TestWinGate: standard against random fails under 70% of the games,
// from either seat order; other line-ups have no win gate.
func TestWinGate(t *testing.T) {
	r := arena.Result{
		Match: arena.Match{Levels: []bot.Level{bot.Random, bot.Standard}, Cap: 150},
		Games: 10, Finished: 10, Wins: map[bot.Level]int{bot.Standard: 6, bot.Random: 4},
		Seats: map[bot.Level]int{bot.Standard: 10, bot.Random: 10},
	}
	var out bytes.Buffer
	if report(&printer{w: &out}, r, 0) || !strings.Contains(out.String(), "standard wins 60% against random") {
		t.Errorf("60%% passed:\n%s", out.String())
	}
	r.Wins = map[bot.Level]int{bot.Standard: 7, bot.Random: 3}
	if !report(&printer{w: &bytes.Buffer{}}, r, 0) {
		t.Error("70% failed")
	}
	if _, ok := winGateOf([]bot.Level{bot.Policy, bot.Random}); ok {
		t.Error("policy against random has a win gate")
	}
	if _, ok := winGateOf([]bot.Level{bot.Standard, bot.Random, bot.Random}); ok {
		t.Error("a 3-player line-up has a win gate")
	}
}

// TestHardGates: hard wins at least 60% of 2-player games against
// standard, and over 30% as one hard among three standard, from any
// seat order; two hard seats among standard have no win gate.
func TestHardGates(t *testing.T) {
	for _, c := range []struct {
		levels []bot.Level
		wins   int
		pass   bool
	}{
		{[]bot.Level{bot.Standard, bot.Hard}, 5, false},
		{[]bot.Level{bot.Hard, bot.Standard}, 6, true},
		{[]bot.Level{bot.Standard, bot.Hard, bot.Standard, bot.Standard}, 3, false},
		{[]bot.Level{bot.Standard, bot.Standard, bot.Standard, bot.Hard}, 4, true},
	} {
		r := arena.Result{
			Match: arena.Match{Levels: c.levels, Cap: 150, Playouts: 150},
			Games: 10, Finished: 10, Wins: map[bot.Level]int{bot.Hard: c.wins, bot.Standard: 10 - c.wins},
			Seats: map[bot.Level]int{bot.Hard: 10, bot.Standard: 10 * (len(c.levels) - 1)},
		}
		var out bytes.Buffer
		if got := report(&printer{w: &out}, r, 0); got != c.pass {
			t.Errorf("%v, hard wins %d of 10: passed %v, want %v:\n%s", c.levels, c.wins, got, c.pass, out.String())
		}
		if !strings.Contains(out.String(), "hard: 150 playouts a move") {
			t.Errorf("%v: no playouts line:\n%s", c.levels, out.String())
		}
	}
	if _, ok := winGateOf([]bot.Level{bot.Hard, bot.Hard, bot.Standard, bot.Standard}); ok {
		t.Error("two hard seats have a win gate")
	}
}

// TestRunBadLevel: an unknown level is an error before any game.
func TestRunBadLevel(t *testing.T) {
	err := run(context.Background(), []string{"-match", "policy,clever"}, &bytes.Buffer{})
	if err == nil || errors.Is(err, errGate) {
		t.Errorf("err %v", err)
	}
}
