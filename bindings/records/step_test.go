package records

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"path/filepath"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// TestStep: steps come by their index in the live steps, with their
// events and view; before Open and out of range is an error.
func TestStep(t *testing.T) {
	r := New(nil, t.TempDir())
	if _, err := r.Step(0); !errors.Is(err, ErrNotOpen) {
		t.Fatalf("before Open: %v", err)
	}
	path := filepath.Join(t.TempDir(), "g"+record.Ext)
	writeRecord(t, path, header("g", time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)), 2, nil)
	if _, err := r.Open(path); err != nil {
		t.Fatal(err)
	}
	s, err := r.Step(3)
	if err != nil {
		t.Fatal(err)
	}
	if s.Index != 3 || s.N != 4 || string(s.Events) != `[]` || string(s.View) != `{"turn":{"number":2,"active":1}}` {
		t.Errorf("step %+v", s)
	}
	for _, i := range []int{-1, 4} {
		if _, err := r.Step(i); !errors.Is(err, ErrNoStep) {
			t.Errorf("step %d: %v", i, err)
		}
	}
}

// TestFullGamePlays: a whole Base Game, recorded as the server records
// it, plays to its end through Open and Step (10.2): every step has its
// events and a view, turns count up, and the last view is the game's end.
func TestFullGamePlays(t *testing.T) {
	path := filepath.Join(t.TempDir(), "full"+record.Ext)
	players := writeFullGame(t, path)
	r := New(nil, t.TempDir())
	sum, err := r.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if sum.End == nil || !sum.End.Finished || len(sum.Turns) < 2 || sum.Steps == 0 {
		t.Fatalf("summary: %d steps, %d turns, end %+v", sum.Steps, len(sum.Turns), sum.End)
	}
	t.Logf("%d steps, %d turns", sum.Steps, len(sum.Turns))
	turn := -1
	var last struct {
		Players []json.RawMessage `json:"players"`
		Turn    struct {
			Number int `json:"number"`
		} `json:"turn"`
		Over    bool  `json:"over"`
		Winners []int `json:"winners"`
	}
	for i := range sum.Steps {
		s, serr := r.Step(i)
		if serr != nil {
			t.Fatalf("step %d: %v", i, serr)
		}
		var events []json.RawMessage
		if jerr := json.Unmarshal(s.Events, &events); jerr != nil {
			t.Fatalf("step %d events: %v", i, jerr)
		}
		if jerr := json.Unmarshal(s.View, &last); jerr != nil {
			t.Fatalf("step %d view: %v", i, jerr)
		}
		if len(last.Players) != players || last.Turn.Number < turn {
			t.Fatalf("step %d: %d players, turn %d after %d", i, len(last.Players), last.Turn.Number, turn)
		}
		turn = last.Turn.Number
	}
	if !last.Over || len(last.Winners) == 0 || turn != sum.Turns[len(sum.Turns)-1].Number {
		t.Errorf("the last view: over %v, winners %v, turn %d", last.Over, last.Winners, turn)
	}
	if _, err := r.Step(sum.Steps); !errors.Is(err, ErrNoStep) {
		t.Errorf("past the end: %v", err)
	}
}

// writeFullGame plays a 2-player Base Game with seeded random moves to
// its end and records it as the server does (internal/server, apply):
// each step with its events and the protocol's judge view. It returns
// the number of players.
func writeFullGame(t *testing.T, path string) int {
	t.Helper()
	base, ok := cards.Find("b2")
	if !ok {
		t.Fatal("no Base Game set")
	}
	const players, maxSteps = 2, 8000
	for seed := uint64(1); seed <= 20; seed++ {
		g, _, err := engine.NewGame(engine.Setup{Seed: seed, Players: players, Sets: []engine.CardSet{base}})
		if err != nil {
			t.Fatal(err)
		}
		h := header("full", time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC))
		w, err := record.Create(path, h)
		if err != nil {
			t.Fatal(err)
		}
		rng := rand.New(rand.NewPCG(seed, 7)) //nolint:gosec // test moves, not security
		for range maxSteps {
			if g.Over {
				break
			}
			p := g.Prompt()
			in := pick(rng, g.Allowed(p.Player))
			if in.Kind == engine.IntentDiscard {
				in.Objects = in.Objects[:p.Count]
			}
			st, serr := g.Step(in)
			if serr != nil {
				t.Fatalf("seed %d: %v", seed, serr)
			}
			view := protocol.ViewJSON(g.View(engine.Viewer{Kind: engine.ViewJudge}))
			if werr := w.Step(st.Intent, st.Events, view, st.Checksum); werr != nil {
				t.Fatal(werr)
			}
		}
		if !g.Over {
			if serr := w.Stop(); serr != nil {
				t.Fatal(serr)
			}
			continue
		}
		winners := make([]int, len(g.Winners))
		for i, p := range g.Winners {
			winners[i] = int(p)
		}
		if cerr := w.Close(record.End{Ended: h.Started.Add(time.Hour), Finished: true, Winners: winners}); cerr != nil {
			t.Fatal(cerr)
		}
		return players
	}
	t.Fatal("no seeded game finished")
	return 0
}

// pick prefers doing something over passing, so the game moves on, as
// the engine's simulations do.
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
