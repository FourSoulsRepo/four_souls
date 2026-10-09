package rulesengine

import (
	"errors"
	"testing"
)

// toAction runs a test game to the active player's open action phase.
func toAction(t *testing.T, players int) *Game {
	t.Helper()
	g := newTestGame(t, players)
	passUntilStep(t, g, StepAction)
	return g
}

func pass(t *testing.T, g *Game) []Event {
	t.Helper()
	ev, err := g.Apply(Intent{Player: g.Prompt().Player, Kind: IntentPass})
	if err != nil {
		t.Fatal(err)
	}
	return ev
}

func resolved(ev []Event) []string {
	var out []string
	for _, e := range ev {
		if e.Kind == EvStackResolved {
			out = append(out, e.Text)
		}
	}
	return out
}

func TestStackIsLastInFirstOut(t *testing.T) {
	g := toAction(t, 3)
	a := g.Turn.Active
	b := g.next(a)
	g.push(StackItem{Kind: StackRoll, Controller: a, Label: "first"})
	g.push(StackItem{Kind: StackRoll, Controller: b, Label: "second"})

	// Priority starts with the controller of the newest item (R-PRIO-03).
	if p := g.Prompt(); p.Player != b {
		t.Fatalf("priority with %d, want %d (R-PRIO-03)", p.Player, b)
	}
	var order []string
	for range 3 {
		order = append(order, resolved(pass(t, g))...)
	}
	if len(order) != 1 || order[0] != "second" {
		t.Fatalf("after one round of passes resolved %v, want [second] (R-STACK-02)", order)
	}
	// After a resolution priority starts with the active player (R-STACK-04).
	if p := g.Prompt(); p.Player != a {
		t.Fatalf("after resolving, priority with %d, want active %d", p.Player, a)
	}
	for range 3 {
		order = append(order, resolved(pass(t, g))...)
	}
	if len(order) != 2 || order[1] != "first" {
		t.Fatalf("resolved %v, want [second first]", order)
	}
	// Empty stack in the action phase: back to the active player (R-TURN-09).
	if !g.inOpenActionPhase() || g.Prompt().Player != a {
		t.Errorf("after the stack emptied: step %d, prompt %+v", g.Turn.Step, g.Prompt())
	}
}

func TestEveryonePassesBeforeResolution(t *testing.T) {
	g := toAction(t, 4)
	g.push(StackItem{Kind: StackRoll, Controller: g.Turn.Active, Label: "x"})
	for i := range 3 {
		if r := resolved(pass(t, g)); len(r) != 0 {
			t.Fatalf("resolved after %d of 4 passes (R-STACK-03)", i+1)
		}
	}
	if r := resolved(pass(t, g)); len(r) != 1 {
		t.Fatal("did not resolve after all 4 players passed")
	}
}

func TestActivePlayerMayPassWithAStack(t *testing.T) {
	g := toAction(t, 2)
	g.push(StackItem{Kind: StackRoll, Controller: g.Turn.Active})
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentPass}); err != nil {
		t.Errorf("with something on the stack the active player may pass: %v", err)
	}
}

func TestNoEndTurnWithAStack(t *testing.T) {
	g := toAction(t, 2)
	g.push(StackItem{Kind: StackRoll, Controller: g.Turn.Active})
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentEndTurn}); err == nil {
		t.Error("ended the turn with something on the stack (R-TURN-06)")
	}
}

func TestPlayLoot(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	card := g.Players[a].Hand[0]
	hand := len(g.Players[a].Hand)
	if _, err := g.Apply(Intent{Player: a, Kind: IntentPlayLoot, Objects: []ObjectID{card}}); err != nil {
		t.Fatal(err)
	}
	if len(g.Stack) != 1 || g.Stack[0].Kind != StackLoot || len(g.Players[a].Hand) != hand-1 {
		t.Fatalf("stack %v, hand %d", g.Stack, len(g.Players[a].Hand))
	}
	if g.Object(g.Stack[0].Source).Zone.Kind != ZoneStack {
		t.Error("a played loot is on the stack (R-MECH-45)")
	}
	// The loot play is used up (R-CARD-08, R-TURN-05).
	next := g.Players[a].Hand[0]
	_, err := g.Apply(Intent{Player: a, Kind: IntentPlayLoot, Objects: []ObjectID{next}})
	var re *RuleError
	if !errors.As(err, &re) || re.Rule != "R-CARD-08" {
		t.Errorf("second loot: %v, want R-CARD-08", err)
	}
	pass(t, g)
	pass(t, g)
	if len(g.Stack) != 0 || len(g.Discards[LootDeck]) != 1 {
		t.Errorf("after resolving: stack %d, loot discard %d (R-CARD-07)", len(g.Stack), len(g.Discards[LootDeck]))
	}
}

func TestOnlyTheActivePlayerHasALootPlay(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	b := g.next(a)
	g.push(StackItem{Kind: StackRoll, Controller: a})
	pass(t, g) // a passes; b gets priority
	if g.Prompt().Player != b {
		t.Fatalf("priority with %d", g.Prompt().Player)
	}
	_, err := g.Apply(Intent{Player: b, Kind: IntentPlayLoot, Objects: []ObjectID{g.Players[b].Hand[0]}})
	if err == nil {
		t.Error("a non-active player played loot without an ability allowing it (R-CARD-08)")
	}
}

func TestRollGoesOnTheStack(t *testing.T) {
	g := toAction(t, 2)
	g.roll(g.Turn.Active, "test")
	if len(g.Stack) != 1 || g.Stack[0].Kind != StackRoll || g.Stack[0].Roll < 1 || g.Stack[0].Roll > 6 {
		t.Fatalf("stack %+v (R-DICE-01, R-DICE-02)", g.Stack)
	}
}
