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

// TestRunBadLevel: an unknown level is an error before any game.
func TestRunBadLevel(t *testing.T) {
	err := run(context.Background(), []string{"-match", "policy,clever"}, &bytes.Buffer{})
	if err == nil || errors.Is(err, errGate) {
		t.Errorf("err %v", err)
	}
}
