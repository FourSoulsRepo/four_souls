package rulesengine

import (
	"errors"
	"testing"
)

// recordGame plays n steps of a seeded game through allowed intents.
func recordGame(t *testing.T, s Setup, n int) (*Game, []StepRecord) {
	t.Helper()
	g, _, err := NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	var steps []StepRecord
	for i := range n {
		if g.Over {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		in := allowed[(i*11)%len(allowed)]
		if in.Kind == IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, err := g.Step(in)
		if err != nil {
			t.Fatal(err)
		}
		steps = append(steps, st)
	}
	return g, steps
}

func TestReplayRebuildsTheSameStates(t *testing.T) {
	s := Setup{Seed: 77, Players: 4, Sets: []CardSet{testSet, abilitySet}}
	g, steps := recordGame(t, s, 300)
	r, err := Replay(s, steps)
	if err != nil {
		t.Fatal(err)
	}
	a, errA := g.Checksum()
	b, errB := r.Checksum()
	if errA != nil || errB != nil {
		t.Fatal(errA, errB)
	}
	if a != b {
		t.Error("the replayed game ends in a different state (A-08)")
	}
}

func TestReplayFindsTheFirstDifference(t *testing.T) {
	s := Setup{Seed: 77, Players: 2, Sets: []CardSet{testSet}}
	_, steps := recordGame(t, s, 50)
	steps[20].Checksum++ // as if a newer engine behaved differently
	_, err := Replay(s, steps)
	if !errors.Is(err, ErrMismatch) {
		t.Fatalf("err = %v, want ErrMismatch", err)
	}
	if err.Error() != "step 20: replay: state differs from the record" {
		t.Errorf("the error names the first differing step: %v", err)
	}
}

func TestSaveAndLoadContinueIdentically(t *testing.T) {
	s := Setup{Seed: 5, Players: 3, Sets: []CardSet{testSet, abilitySet}}
	g, _ := recordGame(t, s, 40)
	data, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	h, err := Load(data, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 60 {
		if g.Over {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		in := allowed[(i*3)%len(allowed)]
		if in.Kind == IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		a, err := g.Step(in)
		if err != nil {
			t.Fatal(err)
		}
		b, err := h.Step(in)
		if err != nil {
			t.Fatal(err)
		}
		if a.Checksum != b.Checksum {
			t.Fatalf("step %d after loading differs", i)
		}
	}
	if _, err := Load(data, testSet); err == nil {
		t.Error("loading with different card sets must fail")
	}
}
