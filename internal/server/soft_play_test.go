package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// softTable starts a 2-player Soft game from a built table, as a loaded
// save: Ann (seat 0) is active in her action phase, Bob is seat 1. The
// table's ids are the game's, so tb.Find names the cards to play.
func softTable(t *testing.T, setup engine.SituationSetup) (*driver, *enginetest.Table) {
	t.Helper()
	return modeTable(t, protocol.RollbackSoft, setup)
}

// modeTable is softTable with another rollback mode.
func modeTable(t *testing.T, mode string, setup engine.SituationSetup) (*driver, *enginetest.Table) {
	t.Helper()
	built := enginetest.NewSetup(t, setup, b2.Set)
	state, err := built.G.Save()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	f := saveFile{
		Format: saveFormat, Engine: engine.Version, Game: "g1", Sets: []string{"b2"},
		Options: protocol.Options{Rollback: mode},
		Seats:   []savedSeat{{Name: "Ann", Token: newToken()}, {Name: "Bob", Token: newToken()}}, State: state,
	}
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "g1"+saveExt), data, 0o600); err != nil {
		t.Fatal(err)
	}
	h := startHubWithSaves(t, dir)
	ann, _ := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeLoad, protocol.Load{Save: "g1"})
	ann.next(protocol.TypeTable, nil)
	bob, _ := arrive(t, h, "Bob", "")
	bob.say(protocol.TypeJoin, protocol.Join{Game: "g1"})
	bob.next(protocol.TypeTable, nil)
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	tb := &driver{t: t, players: []*player{ann, bob}, last: make([]protocol.Update, 2)}
	tb.updates(-1) // the loaded game
	if u := tb.last[0]; u.RollbackMode != mode {
		t.Fatalf("mode %q", u.RollbackMode)
	}
	return tb, built
}

// do sends an intent from its player and returns the step it takes:
// the first step of the update it brings.
func (tb *driver) do(in engine.Intent) int {
	tb.t.Helper()
	k := tb.last[0].Step + 1
	tb.players[in.Player].say(protocol.TypeIntent, protocol.Intent{Intent: in})
	tb.updates(int(in.Player))
	return k
}

// asked is the seat the game waits for; -1: nobody.
func (tb *driver) asked() int {
	for i, u := range tb.last {
		if len(u.Allowed) > 0 {
			return i
		}
	}
	return -1
}

// choose answers the open choose prompt by its option's label.
func (tb *driver) choose(label string) int {
	tb.t.Helper()
	i := tb.asked()
	if i < 0 {
		tb.t.Fatal("nobody is asked")
	}
	w := tb.last[i].View.Waiting
	c := slices.Index(w.Options, label)
	if w.Kind != engine.PromptChoose || c < 0 {
		tb.t.Fatalf("prompt %+v: no option %q", w, label)
	}
	return tb.do(engine.Intent{Player: engine.PlayerID(i), Kind: engine.IntentChoose, Choice: c})
}

// resolve passes until the stack is empty.
func (tb *driver) resolve() {
	tb.t.Helper()
	for range 20 {
		if len(tb.last[0].View.Stack) == 0 {
			return
		}
		i := tb.asked()
		if i < 0 {
			tb.t.Fatal("the stack waits for nobody")
		}
		tb.do(engine.Intent{Player: engine.PlayerID(i), Kind: engine.IntentPass})
	}
	tb.t.Fatalf("the stack never resolved: %+v", tb.last[0].View.Stack)
}

// until plays on, a throw, the first option or a pass, until an update
// brings an event of the kind.
func (tb *driver) until(kind engine.EventKind) {
	tb.t.Helper()
	for range 40 {
		if slices.ContainsFunc(tb.last[0].Events, func(e engine.Event) bool { return e.Kind == kind }) {
			return
		}
		i := tb.asked()
		if i < 0 {
			tb.t.Fatalf("nobody is asked before %s", kind)
		}
		allowed := tb.last[i].Allowed
		in := allowed[0]
		if j := slices.IndexFunc(allowed, func(a engine.Intent) bool { return a.Kind == engine.IntentThrow || a.Kind == engine.IntentChoose }); j >= 0 {
			in = allowed[j]
		} else if j := slices.IndexFunc(allowed, func(a engine.Intent) bool { return a.Kind == engine.IntentPass }); j >= 0 {
			in = allowed[j]
		}
		tb.do(in)
	}
	tb.t.Fatalf("no %s came", kind)
}

// target reports whether step is a Soft target now, and its author.
func (tb *driver) target(step int) (int, bool) {
	i := slices.IndexFunc(tb.last[0].RollbackSteps, func(s protocol.RollbackStep) bool { return s.Step == step })
	if i < 0 {
		return -1, false
	}
	return tb.last[0].RollbackSteps[i].Seat, true
}

// handHas: Ann's hand holds the card.
func (tb *driver) handHas(card engine.CardRef) bool {
	return slices.ContainsFunc(tb.last[0].View.Players[0].Hand, func(c engine.CardView) bool { return c.Card == card })
}

// TestSoftUndoPlayedCard: in Soft a played loot card and an activated
// item are undone by the other player's vote once they resolved, as no
// deck card, die or damage came after (ADR 020, item 4.16; N-12, item
// 24). The table is the one from before the play.
func TestSoftUndoPlayedCard(t *testing.T) {
	for _, c := range []struct {
		name string
		hand []engine.CardRef
		item []engine.CardRef
		in   func(tb *enginetest.Table) engine.Intent
	}{
		{"loot card", []engine.CardRef{"a_penny"}, nil, func(tb *enginetest.Table) engine.Intent {
			return engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "a_penny")}}
		}},
		{"item", nil, []engine.CardRef{"sack_of_pennies"}, func(tb *enginetest.Table) engine.Intent {
			return engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(0, "sack_of_pennies")}}
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			tb, built := softTable(t, engine.SituationSetup{Players: []engine.SituationPlayer{
				{Character: "isaac", Items: c.item, Hand: c.hand}, {Character: "cain"},
			}})
			before := viewJSON(t, tb.last[0])
			k := tb.do(c.in(built))
			tb.resolve()
			if got := tb.last[0].View.Players[0].Cents; got != 1 {
				t.Fatalf("cents %d after the play, want 1", got)
			}
			if author, ok := tb.target(k); !ok || author != 0 {
				t.Fatalf("the play (step %d) is no target of Ann's: from %d, %+v", k, tb.last[0].RollbackFrom, tb.last[0].RollbackSteps)
			}
			tb.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
			tb.rolledBack(k, 1) // two players: Bob's yes is the majority
			if viewJSON(t, tb.last[0]) != before {
				t.Error("the restored table differs from the one before the play")
			}
		})
	}
}

// TestSoftUndoHandToDeck: a card put from the hand on the loot deck is
// no barrier (item 4.22). Bum Friend loots first, which is one; the
// pick of the card to put back is undone, and Ann picks again.
func TestSoftUndoHandToDeck(t *testing.T) {
	tb, built := softTable(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"bum_friend"}, Hand: []engine.CardRef{"a_penny"}}, {Character: "cain"},
	}})
	tb.do(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{built.Find(0, "bum_friend")}})
	tb.until(engine.EvLooted)
	k := tb.choose("a_penny")
	if tb.handHas("a_penny") || !slices.ContainsFunc(tb.last[0].Events, func(e engine.Event) bool { return e.Kind == engine.EvMovedToDeck }) {
		t.Fatalf("A Penny! is not put back: %s", eventsJSON(t, tb.last[0].Events))
	}
	if _, ok := tb.target(k); !ok {
		t.Fatalf("the put back (step %d) is no target: from %d, %+v", k, tb.last[0].RollbackFrom, tb.last[0].RollbackSteps)
	}
	if _, ok := tb.target(k - 1); ok {
		t.Errorf("the loot before it (step %d) is a target", k-1)
	}
	tb.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
	tb.rolledBack(k, 1)
	if !tb.handHas("a_penny") || tb.last[0].View.Waiting.Kind != engine.PromptChoose {
		t.Errorf("Ann is not asked again with A Penny! in hand: %+v", tb.last[0].View.Waiting)
	}
}

// TestSoftPlayedCardBarriers: after a played loot card, each barrier
// stops its rollback: dice, damage, a death, a deck card drawn or shown,
// a hand card seen, and the turn end (ADR 020, items 4 and 4.16).
func TestSoftPlayedCardBarriers(t *testing.T) {
	three := 3
	for _, c := range []struct {
		name string
		kind engine.EventKind
		then func(tb *driver, built *enginetest.Table)
	}{
		{"dice", engine.EvDiceRolled, func(tb *driver, _ *enginetest.Table) {
			tb.do(engine.Intent{Player: 0, Kind: engine.IntentAttack})
		}},
		{"damage", engine.EvDamaged, func(tb *driver, built *enginetest.Table) {
			tb.do(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{built.Find(0, "bomb")}})
		}},
		{"death", engine.EvDied, func(tb *driver, built *enginetest.Table) {
			tb.do(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{built.Find(0, "xiii_death")}})
		}},
		{"monster deck card", engine.EvCardRevealed, func(tb *driver, _ *enginetest.Table) {
			tb.do(engine.Intent{Player: 0, Kind: engine.IntentAttack})
			for range 10 {
				if w := tb.last[0].View.Waiting; w.Kind == engine.PromptChoose && slices.Contains(w.Options, "monster deck") {
					tb.choose("monster deck")
					return
				}
				tb.do(engine.Intent{Player: engine.PlayerID(tb.asked()), Kind: engine.IntentPass})
			}
			tb.t.Fatal("the attack never asks for its target")
		}},
		{"loot deck card", engine.EvLooted, func(tb *driver, built *enginetest.Table) {
			tb.do(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{built.Find(0, "bum_friend")}})
		}},
		{"hand card seen", engine.EvLookedAt, func(tb *driver, built *enginetest.Table) {
			tb.do(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{built.Find(0, "incubus")}})
		}},
		{"turn end", engine.EvTurnEnded, func(tb *driver, _ *enginetest.Table) {
			tb.do(engine.Intent{Player: 0, Kind: engine.IntentEndTurn})
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			tb, built := softTable(t, engine.SituationSetup{
				Players: []engine.SituationPlayer{
					{Character: "isaac", Items: []engine.CardRef{"bum_friend", "incubus"}, Hand: []engine.CardRef{"a_penny", "bomb", "xiii_death"}},
					{Character: "cain", Hand: []engine.CardRef{"a_dime"}},
				},
				Monsters:  []engine.CardRef{"pin"},
				LootPlays: &three,
			})
			k := tb.do(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{built.Find(0, "a_penny")}})
			tb.resolve()
			if _, ok := tb.target(k); !ok {
				t.Fatalf("the played card (step %d) is no target: %+v", k, tb.last[0].RollbackSteps)
			}
			c.then(tb, built)
			tb.until(c.kind)
			if u := tb.last[0]; u.RollbackFrom <= k {
				t.Errorf("after %s the rollback starts at %d, at or before the play (step %d): %s", c.kind, u.RollbackFrom, k, eventsJSON(t, u.Events))
			}
			if _, ok := tb.target(k); ok {
				t.Errorf("after %s the play is still a target", c.kind)
			}
			tb.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
			if code := tb.players[1].failure(); code != protocol.ErrRollback {
				t.Errorf("a rollback past %s: %s", c.kind, code)
			}
		})
	}
}

// TestUndoPlayWithItsTarget: a loot card that asks its target first is
// one move with the answer (17.7). The play is the target; the answer
// is none, so one rollback puts the card back in the hand, in Soft and
// in Godmode, where the targets are listed too.
func TestUndoPlayWithItsTarget(t *testing.T) {
	for _, mode := range []string{protocol.RollbackSoft, protocol.RollbackGod} {
		t.Run(mode, func(t *testing.T) {
			tb, built := modeTable(t, mode, engine.SituationSetup{Players: []engine.SituationPlayer{
				{Character: "isaac", Hand: []engine.CardRef{"bomb"}}, {Character: "cain"},
			}})
			before := viewJSON(t, tb.last[0])
			k := tb.do(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{built.Find(0, "bomb")}})
			if w := tb.last[0].View.Waiting; w.Kind != engine.PromptChoose || w.Player != 0 {
				t.Fatalf("Bomb! asks no target first: %+v", w)
			}
			a := tb.do(engine.Intent{Player: 0, Kind: engine.IntentChoose, Choice: 0})
			if len(tb.last[0].View.Stack) == 0 || tb.handHas("bomb") {
				t.Fatalf("Bomb! is not on the stack: %s", eventsJSON(t, tb.last[0].Events))
			}
			steps := tb.last[0].RollbackSteps
			if len(steps) == 0 || steps[len(steps)-1] != (protocol.RollbackStep{Step: k, Seat: 0}) {
				t.Fatalf("the latest target is not the play (step %d): %+v", k, steps)
			}
			if _, ok := tb.target(a); ok {
				t.Errorf("the answer (step %d) is a target of its own", a)
			}
			if mode == protocol.RollbackSoft {
				tb.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: a})
				if code := tb.players[1].failure(); code != protocol.ErrRollback {
					t.Errorf("a Soft rollback of the answer alone: %s", code)
				}
			}
			tb.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
			tb.rolledBack(k, 1)
			if viewJSON(t, tb.last[0]) != before {
				t.Error("the restored table differs from the one before the play")
			}
		})
	}
}
