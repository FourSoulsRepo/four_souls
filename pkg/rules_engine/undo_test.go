package rulesengine

import (
	"bytes"
	"errors"
	"slices"
	"testing"
)

// undoGame plays a seeded game from a saved base. The first time an
// intent of kind is allowed it is taken; that step is the target, and
// before is the state just before it. The game goes on for after more
// steps, so an undo also undoes what followed.
type undoGame struct {
	g      *Game
	base   []byte
	steps  []StepRecord
	target int
	before []byte
}

func playUndoGame(t *testing.T, s Setup, kind IntentKind, after int) undoGame {
	t.Helper()
	g, _, err := NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	for p := range g.Players {
		g.Players[p].Cents = 30 // so a purchase can pay
	}
	u := undoGame{g: g, target: -1}
	if u.base, err = g.Save(); err != nil {
		t.Fatal(err)
	}
	for i := range 2000 {
		if g.Over || (u.target >= 0 && len(u.steps) > u.target+after) {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		in := allowed[(i*7)%len(allowed)]
		if j := slices.IndexFunc(allowed, func(a Intent) bool { return a.Kind == kind }); u.target < 0 && j >= 0 {
			in = allowed[j]
			u.target = len(u.steps)
			if u.before, err = g.Save(); err != nil {
				t.Fatal(err)
			}
		}
		if in.Kind == IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, err := g.Step(in)
		if err != nil {
			t.Fatal(err)
		}
		u.steps = append(u.steps, st)
	}
	if u.target < 0 {
		t.Fatalf("intent kind %d never allowed", kind)
	}
	return u
}

// checkUndoRedo rolls back to before the target step: the state must be
// the same bytes as before it, the RNG included. Replaying the undone
// steps again must end in the same state as the game (ADR 020, item 7).
func checkUndoRedo(t *testing.T, u undoGame, sets ...CardSet) {
	t.Helper()
	r, err := ReplayFrom(u.base, u.steps[:u.target], sets...)
	if err != nil {
		t.Fatal(err)
	}
	undone, err := r.Save()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(undone, u.before) {
		t.Fatal("the undone game differs from the state before the step")
	}
	r, err = ReplayFrom(undone, u.steps[u.target:], sets...)
	if err != nil {
		t.Fatalf("redo: %v", err)
	}
	a, errA := r.Save()
	b, errB := u.g.Save()
	if errA != nil || errB != nil {
		t.Fatal(errA, errB)
	}
	if !bytes.Equal(a, b) {
		t.Error("undo then redo ends in a different state")
	}
}

// happened finds the first step from i on with an event of kind, or -1.
func happened(steps []StepRecord, i int, kind EventKind) int {
	for ; i < len(steps); i++ {
		if slices.ContainsFunc(steps[i].Events, func(e Event) bool { return e.Kind == kind }) {
			return i
		}
	}
	return -1
}

func TestUndoLootPlay(t *testing.T) {
	sets := []CardSet{testSet, abilitySet}
	u := playUndoGame(t, Setup{Seed: 11, Players: 3, Sets: sets}, IntentPlayLoot, 30)
	if happened(u.steps, u.target, EvLootPlayed) != u.target {
		t.Fatal("the target step played no loot card")
	}
	checkUndoRedo(t, u, sets...)
}

func TestUndoAttack(t *testing.T) {
	sets := []CardSet{testSet, abilitySet}
	u := playUndoGame(t, Setup{Seed: 12, Players: 2, Sets: sets}, IntentAttack, 30)
	if happened(u.steps, u.target, EvAttackDeclared) != u.target || happened(u.steps, u.target, EvDiceRolled) < 0 {
		t.Fatal("the target step declared no attack, or no attack roll followed")
	}
	checkUndoRedo(t, u, sets...)
}

func TestUndoPurchase(t *testing.T) {
	sets := []CardSet{testSet, abilitySet}
	u := playUndoGame(t, Setup{Seed: 13, Players: 4, Sets: sets}, IntentPurchase, 30)
	if happened(u.steps, u.target, EvPurchaseDeclared) != u.target || happened(u.steps, u.target, EvPurchased) < 0 {
		t.Fatal("the target step declared no purchase, or nothing was bought")
	}
	checkUndoRedo(t, u, sets...)
}

// TestUndoRollRepeats: roll, roll back, roll again gives the same value.
// The roll is an attack roll well into the game. The restored RNG
// repeats it; nothing reseeds (N-12, item 7; ADR 020, item 11).
func TestUndoRollRepeats(t *testing.T) {
	sets := []CardSet{testSet, abilitySet}
	u := playUndoGame(t, Setup{Seed: 14, Players: 2, Sets: sets}, IntentAttack, 30)
	k := happened(u.steps, u.target, EvDiceRolled)
	if k < 0 {
		t.Fatal("no attack roll followed the attack")
	}
	rolled := func(ev []Event) Event {
		return ev[slices.IndexFunc(ev, func(e Event) bool { return e.Kind == EvDiceRolled })]
	}
	first := rolled(u.steps[k].Events)
	r, err := ReplayFrom(u.base, u.steps[:k], sets...)
	if err != nil {
		t.Fatal(err)
	}
	ev, err := r.Apply(u.steps[k].Intent)
	if err != nil {
		t.Fatal(err)
	}
	if again := rolled(ev); again.Amount != first.Amount || !slices.Equal(again.Dice, first.Dice) {
		t.Errorf("rolled %d, then %d after the rollback", first.Amount, again.Amount)
	}
}

// TestReplayFromAnyBase: a base saved mid-game rebuilds every later
// state; the checksum of each step is checked.
func TestReplayFromAnyBase(t *testing.T) {
	sets := []CardSet{testSet, abilitySet}
	s := Setup{Seed: 15, Players: 3, Sets: sets}
	g, steps := recordGame(t, s, 40)
	base, saveErr := g.Save()
	if saveErr != nil {
		t.Fatal(saveErr)
	}
	for i := range 60 {
		if g.Over {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		in := allowed[(i*5)%len(allowed)]
		if in.Kind == IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, err := g.Step(in)
		if err != nil {
			t.Fatal(err)
		}
		steps = append(steps, st)
	}
	later := steps[40:]
	r, err := ReplayFrom(base, later, sets...)
	if err != nil {
		t.Fatal(err)
	}
	a, errA := r.Checksum()
	b, errB := g.Checksum()
	if errA != nil || errB != nil {
		t.Fatal(errA, errB)
	}
	if a != b {
		t.Error("the game rebuilt from a mid-game base ends in a different state")
	}
	later[10].Checksum++ // as if the step had gone another way
	if _, err := ReplayFrom(base, later, sets...); !errors.Is(err, ErrMismatch) || err.Error() != "step 10: replay: state differs from the record" {
		t.Errorf("err = %v, want ErrMismatch at step 10", err)
	}
	if _, err := ReplayFrom(base, later, testSet); err == nil {
		t.Error("a base loaded with different card sets must fail")
	}
}
