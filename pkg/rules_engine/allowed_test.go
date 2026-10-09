package rulesengine

import (
	"errors"
	"testing"
)

// playRandom drives a game with intents picked from Allowed, checking on
// a clone that every allowed intent is accepted.
func TestAllowedIntentsAreAccepted(t *testing.T) {
	g, _, err := NewGame(Setup{Seed: 21, Players: 3, Sets: []CardSet{testSet, abilitySet}})
	if err != nil {
		t.Fatal(err)
	}
	giveItem(g, 0, "coin_bag")
	giveItem(g, 1, "razor")
	for step := range 400 {
		if g.Over {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		if len(allowed) == 0 {
			t.Fatalf("step %d: no allowed intent for the waiting player, prompt %+v", step, p)
		}
		for _, in := range allowed {
			c, err := g.Clone()
			if err != nil {
				t.Fatal(err)
			}
			if in.Kind == IntentDiscard {
				in.Objects = in.Objects[:p.Count]
			}
			if _, err := c.Apply(in); err != nil {
				t.Fatalf("step %d: allowed intent %+v refused: %v", step, in, err)
			}
		}
		// Others have nothing allowed while one player is asked.
		for _, pl := range g.Players {
			if pl.ID != p.Player && len(g.Allowed(pl.ID)) > 0 {
				t.Fatalf("player %d has allowed intents while player %d is asked", pl.ID, p.Player)
			}
		}
		in := allowed[(step*7)%len(allowed)]
		if in.Kind == IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRefusalsCiteARule(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	b := g.next(a)
	cases := []Intent{
		{Player: b, Kind: IntentPass},
		{Player: b, Kind: IntentEndTurn},
		{Player: a, Kind: IntentPass},
		{Player: a, Kind: IntentDiscard},
		{Player: a, Kind: IntentChoose},
		{Player: a, Kind: IntentPlayLoot, Objects: []ObjectID{9999}},
	}
	for _, in := range cases {
		_, err := g.Apply(in)
		var re *RuleError
		if !errors.As(err, &re) || re.Rule == "" || re.Reason == "" {
			t.Errorf("%+v: refusal %v must carry a rule ID and a reason", in, err)
		}
	}
}

func TestCanAct(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	b := g.next(a)
	if !g.CanAct(a) {
		t.Error("the active player in the action phase can act")
	}
	g.push(StackItem{Kind: StackRoll, Controller: a})
	pass(t, g)
	if g.Prompt().Player != b {
		t.Fatal("b should have priority")
	}
	// b has no loot play and no abilities: only pass.
	if g.CanAct(b) {
		t.Errorf("b can only pass; allowed %+v", g.Allowed(b))
	}
}
