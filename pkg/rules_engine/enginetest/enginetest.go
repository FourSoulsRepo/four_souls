// Package enginetest helps write card tests: build a table, act, answer
// prompts by label and check the result.
//
//	tb := enginetest.New(t, b2.Set, enginetest.Seat("isaac", "the_d6"), enginetest.Seat("cain"))
//	tb.Activate(0, "isaac", 0)
//	if tb.G.Turn.LootPlays != 2 { … }
package enginetest

import (
	"errors"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Table is a game under test.
type Table struct {
	T testing.TB
	G *engine.Game
}

// Seat is a player with a character and items in play.
func Seat(character engine.CardRef, items ...engine.CardRef) engine.SituationPlayer {
	return engine.SituationPlayer{Character: character, Items: items}
}

// New builds a table where player 0 is active, in the open action phase.
func New(t testing.TB, set engine.CardSet, seats ...engine.SituationPlayer) *Table {
	t.Helper()
	return NewSetup(t, engine.SituationSetup{Players: seats}, set)
}

// NewSetup builds a table from a full setup.
func NewSetup(t testing.TB, setup engine.SituationSetup, sets ...engine.CardSet) *Table {
	t.Helper()
	g, err := engine.BuildTable(setup, sets...)
	if err != nil {
		t.Fatal(err)
	}
	return &Table{T: t, G: g}
}

// Find returns the first object p controls (character or in play) or
// holds in hand with the given card.
func (tb *Table) Find(p int, card engine.CardRef) engine.ObjectID {
	tb.T.Helper()
	pl := tb.G.Players[p]
	ids := append([]engine.ObjectID{pl.Character}, pl.InPlay...)
	ids = append(ids, pl.Hand...)
	for _, id := range ids {
		if tb.G.Object(id).Card == card {
			return id
		}
	}
	tb.T.Fatalf("player %d has no %s", p, card)
	return 0
}

// Do applies an intent and settles with the given choices; it fails the
// test if the engine refuses. It returns every event on the way.
func (tb *Table) Do(in engine.Intent, choices ...string) []engine.Event {
	tb.T.Helper()
	ev, err := tb.G.Apply(in)
	if err != nil {
		tb.T.Fatalf("%+v refused: %v", in, err)
	}
	return append(ev, tb.Settle(choices...)...)
}

// Settle answers prompts by label and passes until the stack is empty.
func (tb *Table) Settle(choices ...string) []engine.Event {
	tb.T.Helper()
	ev, err := tb.G.Settle(choices...)
	if err != nil {
		tb.T.Fatal(err)
	}
	return ev
}

// Activate uses ability i of a card p controls, answering the prompts.
func (tb *Table) Activate(p int, card engine.CardRef, i int, choices ...string) []engine.Event {
	tb.T.Helper()
	return tb.Do(engine.Intent{Player: engine.PlayerID(p), Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(p, card)}, Choice: i}, choices...)
}

// Refused returns the rule that refuses activating ability i, or fails
// if it is allowed.
func (tb *Table) Refused(p int, card engine.CardRef, i int) string {
	tb.T.Helper()
	_, err := tb.G.Apply(engine.Intent{Player: engine.PlayerID(p), Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(p, card)}, Choice: i})
	var re *engine.RuleError
	if !errors.As(err, &re) {
		tb.T.Fatalf("activating %s ability %d: want a refusal, got %v", card, i, err)
	}
	return re.Rule
}

// EndTurn ends the active player's turn and runs the game up to the next
// player's action phase, answering prompts with choices on the way.
func (tb *Table) EndTurn(choices ...string) {
	tb.T.Helper()
	active := tb.G.Turn.Active
	if _, err := tb.G.Apply(engine.Intent{Player: active, Kind: engine.IntentEndTurn}); err != nil {
		tb.T.Fatal(err)
	}
	for range 1000 {
		if tb.G.Turn.Active != active && tb.G.Turn.Step == engine.StepAction && tb.G.Prompt().Kind == engine.PromptPriority && len(tb.G.Stack) == 0 {
			return
		}
		w := tb.G.Prompt()
		var in engine.Intent
		switch w.Kind {
		case engine.PromptChoose:
			if len(choices) == 0 {
				tb.T.Fatalf("a choice is needed, options %v", w.Options)
			}
			in = engine.Intent{Player: w.Player, Kind: engine.IntentChoose, Choice: tb.option(w, choices[0])}
			choices = choices[1:]
		case engine.PromptDiscard:
			in = engine.Intent{Player: w.Player, Kind: engine.IntentDiscard, Objects: tb.G.Players[w.Player].Hand[:w.Count]}
		case engine.PromptPriority:
			in = engine.Intent{Player: w.Player, Kind: engine.IntentPass}
		default:
			tb.T.Fatalf("unexpected prompt %+v", w)
		}
		if _, err := tb.G.Apply(in); err != nil {
			tb.T.Fatal(err)
		}
	}
	tb.T.Fatal("the next turn never came")
}

func (tb *Table) option(w engine.Prompt, label string) int {
	tb.T.Helper()
	for i, o := range w.Options {
		if o == label {
			return i
		}
	}
	tb.T.Fatalf("%q is not an option of %v", label, w.Options)
	return 0
}

// Options returns the labels of the open choose prompt.
func (tb *Table) Options() []string {
	return tb.G.Prompt().Options
}

// Attack makes the active player attack the monster with the given card,
// with the given dice results, and settles.
func (tb *Table) Attack(monster engine.CardRef, rolls ...int) []engine.Event {
	tb.T.Helper()
	tb.G.ForceRolls(rolls...)
	return tb.Do(engine.Intent{Player: tb.G.Turn.Active, Kind: engine.IntentAttack}, string(monster))
}

// Has reports whether events include one of the kind for player p.
func Has(events []engine.Event, kind engine.EventKind, p engine.PlayerID) bool {
	for _, e := range events {
		if e.Kind == kind && e.Player == p {
			return true
		}
	}
	return false
}

// Pass passes priority for p without settling, e.g. so the next player
// may act on the active player's turn.
func (tb *Table) Pass(p int) {
	tb.T.Helper()
	if _, err := tb.G.Apply(engine.Intent{Player: engine.PlayerID(p), Kind: engine.IntentPass}); err != nil {
		tb.T.Fatal(err)
	}
}

// Start activates ability i without settling: the ability waits on the
// stack, so other players may respond.
func (tb *Table) Start(p int, card engine.CardRef, i int) {
	tb.T.Helper()
	in := engine.Intent{Player: engine.PlayerID(p), Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(p, card)}, Choice: i}
	if _, err := tb.G.Apply(in); err != nil {
		tb.T.Fatalf("%+v refused: %v", in, err)
	}
}
