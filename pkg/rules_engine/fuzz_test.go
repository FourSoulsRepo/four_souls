package rulesengine

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// fuzzSets are the cards fuzzing plays with; the real sets join in step 5.
var fuzzSets = []CardSet{testSet, abilitySet}

// maxFuzzSteps caps one fuzz game. Every step hashes the whole state, so
// very long inputs would slow fuzzing down without finding more.
const maxFuzzSteps = 2000

// playFuzz plays a game from fuzz input: seed, player count and a stream
// of choices among allowed intents. focus cards start in every player's
// play area. It checks the invariants after every step and returns the
// steps for a replay check.
func playFuzz(t *testing.T, seed uint64, players uint8, choices []byte, focus []CardRef) {
	t.Helper()
	s := Setup{Seed: seed, Players: 2 + int(players%3), Sets: fuzzSets, BonusSouls: true}
	g, _, err := NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, card := range focus {
		for p := range g.Players {
			giveItem(g, PlayerID(p), card)
		}
	}
	if err = g.CheckInvariants(); err != nil {
		t.Fatalf("after setup: %v", err)
	}
	base, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	var steps []StepRecord
	choices = choices[:min(len(choices), maxFuzzSteps)]
	for i, c := range choices {
		if g.Over {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		if len(allowed) == 0 {
			t.Fatalf("step %d: nothing allowed for the waiting player, prompt %+v", i, p)
		}
		in := allowed[int(c)%len(allowed)]
		if in.Kind == IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, err := g.Step(in)
		if err != nil {
			t.Fatalf("step %d: allowed intent %+v refused: %v", i, in, err)
		}
		steps = append(steps, st)
		if err := g.CheckInvariants(); err != nil {
			t.Fatalf("step %d after %+v: %v", i, in, err)
		}
		checkNoLeaks(t, g)
	}
	if len(focus) == 0 { // a replay starts without the focus items
		if _, err := Replay(s, steps); err != nil {
			t.Fatalf("replay: %v (A-08)", err)
		}
	}
	checkUndoRedoFuzz(t, g, base, steps, undoStep(choices, len(steps)))
}

// undoStep picks the step to roll back to before, 0 to n, from the input.
func undoStep(choices []byte, n int) int {
	if len(choices) == 0 {
		return 0
	}
	return int(choices[len(choices)-1]) * n / 255
}

// checkUndoRedoFuzz rolls the game back to before step k and plays the
// undone steps again. Each step's checksum covers the whole state, the
// RNG included, so the undo is exact; the redo must end in the same
// bytes as the game (ADR 020, item 7).
func checkUndoRedoFuzz(t *testing.T, g *Game, base []byte, steps []StepRecord, k int) {
	t.Helper()
	r, err := ReplayFrom(base, steps[:k], fuzzSets...)
	if err != nil {
		t.Fatalf("undo to before step %d: %v", k, err)
	}
	undone, err := r.Save()
	if err != nil {
		t.Fatal(err)
	}
	if r, err = ReplayFrom(undone, steps[k:], fuzzSets...); err != nil {
		t.Fatalf("redo from step %d: %v", k, err)
	}
	a, errA := r.Save()
	b, errB := g.Save()
	if errA != nil || errB != nil {
		t.Fatal(errA, errB)
	}
	if !bytes.Equal(a, b) {
		t.Fatalf("undo to before step %d, then redo, ends in a different state", k)
	}
}

func checkNoLeaks(t *testing.T, g *Game) {
	t.Helper()
	// No view names a card hidden in a deck (A-16). The top treasure for
	// a peeking player is allowed; viewIDs leaves it out.
	inDeck := map[ObjectID]bool{}
	for d := range deckCount {
		for _, id := range g.Decks[d] {
			inDeck[id] = true
		}
	}
	viewers := []Viewer{{Kind: ViewSpectator}}
	for _, pl := range g.Players {
		viewers = append(viewers, Viewer{Kind: ViewPlayer, Player: pl.ID})
	}
	for _, v := range viewers {
		gv := g.View(v)
		for id := range viewIDs(gv) {
			if inDeck[id] {
				t.Fatalf("viewer %+v sees %s, hidden in a deck (A-16)", v, g.Object(id).Card)
			}
		}
		for _, it := range gv.Stack {
			if it.Card != "" && inDeck[it.Source] {
				t.Fatalf("viewer %+v sees stack item %s from a deck (A-16)", v, it.Card)
			}
		}
	}
	for _, pl := range g.Players {
		ids := viewIDs(g.View(Viewer{Kind: ViewPlayer, Player: pl.ID}))
		for _, other := range g.Players {
			if other.ID == pl.ID {
				continue
			}
			for _, h := range other.Hand {
				if ids[h] {
					t.Fatalf("player %d sees a hand card of player %d (A-07)", pl.ID, other.ID)
				}
			}
		}
	}
}

// FuzzAllCards plays random games with every card (TS-05, "all cards").
// CI runs only the seed corpus; real fuzzing runs locally:
//
//	go test -run '^$' -fuzz=FuzzAllCards -fuzztime=10m -fuzzminimizetime=5s .
func FuzzAllCards(f *testing.F) {
	f.Add(uint64(1), uint8(0), []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	f.Add(uint64(42), uint8(2), []byte(strings.Repeat("\x01\x00\x03\x02", 60)))
	f.Add(uint64(7), uint8(1), []byte(strings.Repeat("\x00", 200)))
	f.Fuzz(func(t *testing.T, seed uint64, players uint8, choices []byte) {
		playFuzz(t, seed, players, choices, nil)
	})
}

// FuzzFocused keeps the cards named in FOUR_SOULS_FOCUS (comma-separated)
// in every player's play area (TS-05, "focused" mode for new cards):
//
//	FOUR_SOULS_FOCUS=razor,d6 go test -run '^$' -fuzz=FuzzFocused -fuzztime=5m .
func FuzzFocused(f *testing.F) {
	focus := []CardRef{"razor", "coin_bag", "d6"}
	if env := os.Getenv("FOUR_SOULS_FOCUS"); env != "" {
		focus = nil
		for _, c := range strings.Split(env, ",") {
			focus = append(focus, CardRef(strings.TrimSpace(c)))
		}
	}
	f.Add(uint64(3), uint8(1), []byte(strings.Repeat("\x02\x05\x01", 80)))
	f.Add(uint64(9), uint8(0), []byte{1, 1, 1, 0, 0, 2, 3, 4})
	f.Fuzz(func(t *testing.T, seed uint64, players uint8, choices []byte) {
		playFuzz(t, seed, players, choices, focus)
	})
}
