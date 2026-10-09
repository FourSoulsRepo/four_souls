package rulesengine

import (
	"testing"
)

// Fake items with replacement effects, controlled by their owner.
var (
	// plusOne: "If you would gain ¢, gain that much +1¢ instead" (R-ABIL-29).
	plusOne = Replacement{
		Text: "+1¢",
		When: func(g *Game, self ObjectID, a Action) bool {
			return a.Kind == ActGainCents && a.Player == g.Object(self).Controller
		},
		Do: func(_ *Game, _ ObjectID, a Action) []Action {
			a.Amount++
			return []Action{a}
		},
	}
	// double: "If you would gain ¢, gain double instead".
	double = Replacement{
		Text: "double ¢",
		When: func(g *Game, self ObjectID, a Action) bool {
			return a.Kind == ActGainCents && a.Player == g.Object(self).Controller
		},
		Do: func(_ *Game, _ ObjectID, a Action) []Action {
			a.Amount *= 2
			return []Action{a}
		},
	}
	// keepCents: "Prevent losing ¢" (prevent replaces with nothing, R-ABIL-30).
	keepCents = Replacement{
		Text: "prevent losing ¢",
		When: func(g *Game, self ObjectID, a Action) bool {
			return a.Kind == ActLoseCents && a.Player == g.Object(self).Controller
		},
		Do: func(*Game, ObjectID, Action) []Action { return nil },
	}
	// lootDouble: "If you would loot, loot double instead".
	lootDouble = Replacement{
		Text: "loot double",
		When: func(g *Game, self ObjectID, a Action) bool {
			return a.Kind == ActLoot && a.Player == g.Object(self).Controller
		},
		Do: func(_ *Game, _ ObjectID, a Action) []Action {
			a.Amount *= 2
			return []Action{a}
		},
	}
)

var replacementSet = CardSet{Name: "replacements", Cards: []CardDef{
	{Ref: "plus_one", Kind: TreasureCard, Outside: true, Replacements: []Replacement{plusOne}},
	{Ref: "doubler", Kind: TreasureCard, Outside: true, Replacements: []Replacement{double}},
	{Ref: "wallet", Kind: TreasureCard, Outside: true, Replacements: []Replacement{keepCents}},
	{Ref: "two_clubs", Kind: TreasureCard, Outside: true, Replacements: []Replacement{lootDouble}},
}}

func newReplacementGame(t *testing.T) *Game {
	t.Helper()
	g, _, err := NewGame(Setup{Seed: 9, Players: 2, Sets: []CardSet{testSet, replacementSet}})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// giveItem puts a card into play under p's control.
func giveItem(g *Game, p PlayerID, card CardRef) ObjectID {
	id := g.newObject(card, Zone{Kind: ZoneInPlay}, p)
	o := g.Object(id)
	o.Role, o.Charged = RoleItem, true
	g.Players[p].InPlay = append(g.Players[p].InPlay, id)
	return id
}

// runAction queues an action and drains the queue as if a stack item were
// resolving; the prompt stays as it was unless a choice is asked.
func runAction(g *Game, a Action) []Event {
	saved := g.Waiting
	g.Waiting = Prompt{}
	g.enqueue(a)
	for len(g.Queue) > 0 && g.Waiting.Kind == PromptNone {
		g.processQueue()
	}
	if g.Waiting.Kind == PromptNone {
		g.Waiting = saved
	}
	return g.takeEvents()
}

func TestReplacementChangesAction(t *testing.T) {
	g := newReplacementGame(t)
	p := PlayerID(0)
	giveItem(g, p, "plus_one")
	before := g.Players[p].Cents
	runAction(g, Action{Kind: ActGainCents, Player: p, Amount: 3})
	if got := g.Players[p].Cents - before; got != 4 {
		t.Errorf("gained %d, want 3+1 (R-ABIL-29)", got)
	}
	// It does not touch another player's gain.
	before1 := g.Players[1].Cents
	runAction(g, Action{Kind: ActGainCents, Player: 1, Amount: 3})
	if got := g.Players[1].Cents - before1; got != 3 {
		t.Errorf("player 1 gained %d, want 3", got)
	}
}

func TestPreventReplacesWithNothing(t *testing.T) {
	g := newReplacementGame(t)
	giveItem(g, 0, "wallet")
	before := g.Players[0].Cents
	ev := runAction(g, Action{Kind: ActLoseCents, Player: 0, Amount: 2})
	if g.Players[0].Cents != before {
		t.Errorf("cents %d -> %d; the loss was prevented (R-ABIL-30)", before, g.Players[0].Cents)
	}
	for _, e := range ev {
		if e.Kind == EvLostCents {
			t.Error("a prevented action must not emit its event")
		}
	}
}

func TestReplacementAppliesOnce(t *testing.T) {
	g := newReplacementGame(t)
	giveItem(g, 0, "plus_one")
	before := g.Players[0].Cents
	runAction(g, Action{Kind: ActGainCents, Player: 0, Amount: 1})
	if got := g.Players[0].Cents - before; got != 2 {
		t.Errorf("gained %d; a replacement applies once per action, want 2 (R-ABIL-32)", got)
	}
}

func TestPlayerOrdersSeveralReplacements(t *testing.T) {
	for _, tc := range []struct {
		first string
		want  int
	}{
		{"+1¢", (3 + 1) * 2},
		{"double ¢", 3*2 + 1},
	} {
		g := newReplacementGame(t)
		giveItem(g, 0, "plus_one")
		giveItem(g, 0, "doubler")
		before := g.Players[0].Cents
		runAction(g, Action{Kind: ActGainCents, Player: 0, Amount: 3})
		pr := g.Prompt()
		if pr.Kind != PromptChoose || pr.Purpose != ChooseReplacement || pr.Player != 0 || len(pr.Options) != 2 {
			t.Fatalf("prompt %+v; the affected player orders the replacements (R-ABIL-33)", pr)
		}
		choice := 0
		for i, o := range pr.Options {
			if o == tc.first {
				choice = i
			}
		}
		if _, err := g.Apply(Intent{Player: 1, Kind: IntentChoose, Choice: choice}); err == nil {
			t.Fatal("another player chose the order")
		}
		if _, err := g.Apply(Intent{Player: 0, Kind: IntentChoose, Choice: 5}); err == nil {
			t.Fatal("an out-of-range choice was accepted")
		}
		if _, err := g.Apply(Intent{Player: 0, Kind: IntentChoose, Choice: choice}); err != nil {
			t.Fatal(err)
		}
		if got := g.Players[0].Cents - before; got != tc.want {
			t.Errorf("%s first: gained %d, want %d", tc.first, got, tc.want)
		}
	}
}

func TestLootStepGoesThroughReplacements(t *testing.T) {
	g, _, err := NewGame(Setup{Seed: 9, Players: 2, Sets: []CardSet{testSet, replacementSet}})
	if err != nil {
		t.Fatal(err)
	}
	// Give the next player Two-of-Clubs style doubling before their turn.
	nextP := g.next(g.Turn.Active)
	giveItem(g, nextP, "two_clubs")
	passUntilStep(t, g, StepAction)
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	hand := len(g.Players[nextP].Hand)
	for g.Turn.Active != nextP || g.Turn.Step != StepAction {
		p := g.Prompt()
		in := Intent{Player: p.Player, Kind: IntentPass}
		if p.Kind == PromptDiscard {
			in = Intent{Player: p.Player, Kind: IntentDiscard, Objects: g.Players[p.Player].Hand[:p.Count]}
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatal(err)
		}
	}
	if got := len(g.Players[nextP].Hand) - hand; got != 2 {
		t.Errorf("loot step drew %d, want 2 with loot doubling", got)
	}
}
