package rulesengine

import (
	"fmt"
	"strings"
	"testing"
)

// deckAttackState is what a deck attack leaves once its revealed card
// is placed and the first attack roll waits: the slots, the attack and
// the turn's attacks.
func deckAttackState(g *Game) string {
	out := ""
	for i, s := range g.Monsters {
		out += fmt.Sprintf("slot %d:", i+1)
		for _, id := range s.Cards {
			o := g.Object(id)
			out += fmt.Sprintf(" %s(zone %d, role %d)", o.Card, o.Zone.Kind, o.Role)
		}
		out += "\n"
	}
	a := g.Attack
	top, _ := g.Monsters[0].TopOf()
	out += fmt.Sprintf("attack on %v, started %v, deck %v, target %s, target is slot 1's top %v, on a monster %v\n",
		a.On, a.Started, a.Deck, g.Object(a.Target).Card, a.Target == top, g.AttackFor(NeedsAttackOnMonster, g.Turn.Active))
	out += fmt.Sprintf("attacks %d, left %d, deck attacks %d, must attack deck %d\n",
		g.Turn.Attacks, g.attacksLeft(), g.Turn.DeckAttacks, g.Turn.MustAttackDeck)
	it := g.Stack[len(g.Stack)-1]
	out += fmt.Sprintf("stack %d, top: roll %v, attack %v, by %d, at the target %v",
		len(g.Stack), it.Kind == StackRoll, it.Attack, it.Controller, it.Target.Object == a.Target)
	return out
}

// A stacked attack roll with from_deck builds what a real deck attack
// builds (24.33): the revealed card covers the slot's card (R-ATK-08,
// R-ZONE-10), the attack is on it and counts as one on a monster
// (R-ATK-05, Q-24), and one attack is spent (R-TURN-07). A dummy in
// slot 2 is another copy: the roll attacks the revealed one.
func TestStackedDeckAttackRollMatchesADeckAttack(t *testing.T) {
	setup := SituationSetup{
		Players:  []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}},
		Monsters: []CardRef{"gaper", "dummy"},
	}
	g, err := BuildTable(setup, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	g.PutOnDeck(MonsterDeck, "dummy")
	attack(t, g)
	choose(t, g, "ask.deck.monster")
	choose(t, g, g.Plain(SlotOption(1, "gaper")))
	if w := g.Prompt(); w.Kind != PromptThrow {
		t.Fatalf("prompt %+v; the attack roll waits for its throw", w)
	}
	want := deckAttackState(g)

	setup.Stack = []SituationStack{{Kind: "attack_roll", DiceRoll: 6, Card: "dummy", Owner: 0, FromDeck: true, Slot: 1}}
	s, err := BuildTable(setup, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	if got := deckAttackState(s); got != want {
		t.Errorf("stacked deck attack roll:\n%s\nwant, as a real deck attack:\n%s", got, want)
	}
	if s.Stack[0].Roll != 6 {
		t.Errorf("the roll is %d; want the setup's 6", s.Stack[0].Roll)
	}
}

// A second roll of the same attack names the revealed card without
// from_deck: it attacks the same copy, and no second attack is spent.
func TestStackedDeckAttackSecondRoll(t *testing.T) {
	g, err := BuildTable(SituationSetup{
		Players:  []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}},
		Monsters: []CardRef{"dummy", "gaper"},
		Stack: []SituationStack{
			{Kind: "attack_roll", DiceRoll: 1, Card: "dummy", Owner: 0, FromDeck: true, Slot: 2},
			{Kind: "attack_roll", DiceRoll: 2, Card: "dummy", Owner: 0},
		},
	}, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	top, _ := g.Monsters[1].TopOf()
	if g.Attack.Target != top || !g.Attack.Deck {
		t.Errorf("attack %+v; want the deck attack on slot 2's revealed dummy (%d)", g.Attack, top)
	}
	for i, it := range g.Stack {
		if it.Target.Object != top {
			t.Errorf("roll %d attacks %d; want the revealed dummy %d", i+1, it.Target.Object, top)
		}
	}
	if g.attacksLeft() != 0 {
		t.Errorf("attacks left %d; want 0: one attack, spent by the deck attack (R-TURN-07)", g.attacksLeft())
	}
}

// A wrong deck attack roll is an error that names what is wrong.
func TestStackedDeckAttackRollErrors(t *testing.T) {
	stone := CardSet{Name: "stone", Cards: []CardDef{{Ref: "stone", Kind: MonsterCard, HP: 1, DC: 1, Copies: 1, Unattackable: true}}}
	const prefix = "situation: stack item "
	for _, c := range []struct {
		stack []SituationStack
		want  string
	}{
		{
			[]SituationStack{{Kind: "attack_roll", DiceRoll: 3, Card: "dummy", FromDeck: true}},
			`1: from_deck needs slot: the monster slot "dummy" covers, from 1 (R-ATK-08)`,
		},
		{
			[]SituationStack{{Kind: "attack_roll", DiceRoll: 3, Card: "gaper", Slot: 1}},
			`1: slot is the monster slot a card revealed from the monster deck covers; set "from_deck": true, or leave slot out`,
		},
		{
			[]SituationStack{{Kind: "play", Card: "penny", FromDeck: true, Slot: 1}},
			`1: kind "play" has no from_deck or slot; they are for an attack_roll`,
		},
		{
			[]SituationStack{{Kind: "attack_roll", DiceRoll: 3, Card: "dummy", FromDeck: true, Slot: 3}},
			"1: from_deck: slot 3 is no monster slot, options 1 to 2 (R-ATK-08)",
		},
		{
			[]SituationStack{{Kind: "attack_roll", DiceRoll: 3, Card: "penny", FromDeck: true, Slot: 1}},
			`1: from_deck: "penny" is no card of the monster deck, so no deck attack reveals it (R-ATK-08)`,
		},
		{
			[]SituationStack{{Kind: "attack_roll", DiceRoll: 3, Card: "ambush_event", FromDeck: true, Slot: 1}},
			`1: from_deck: "ambush_event" is an event; revealed, it triggers and the attack is over (R-ATK-09), so no attack roll waits`,
		},
		{
			[]SituationStack{{Kind: "attack_roll", DiceRoll: 3, Card: "stone", FromDeck: true, Slot: 1}},
			`1: from_deck: "stone" can't be attacked; revealed, it stays and the attack is cancelled (R-ATK-10), so no attack roll waits`,
		},
		{[]SituationStack{
			{Kind: "attack_roll", DiceRoll: 3, Card: "dummy", FromDeck: true, Slot: 1},
			{Kind: "attack_roll", DiceRoll: 3, Card: "dummy", FromDeck: true, Slot: 2},
		}, `2: from_deck: an attack on the monster deck reveals one card, on the first attack roll (R-ATK-08); a later roll of the attack names "dummy" without from_deck and slot`},
		{[]SituationStack{
			{Kind: "attack_roll", DiceRoll: 3, Card: "dummy", FromDeck: true, Slot: 1},
			{Kind: "attack_roll", DiceRoll: 3, Card: "boss"},
		}, "2: every attack roll must attack the same monster"},
	} {
		_, err := BuildTable(SituationSetup{
			Players:  []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}},
			Monsters: []CardRef{"gaper", "boss"},
			Stack:    c.stack,
		}, testSet, abilitySet, stone)
		if want := prefix + c.want; err == nil || err.Error() != want {
			t.Errorf("error %v;\nwant %q", err, want)
		}
	}
	// With no attack left it is refused as any attack roll (R-TURN-07).
	none := 0
	_, err := BuildTable(SituationSetup{
		Players:  []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}},
		Monsters: []CardRef{"gaper"},
		Attacks:  &none,
		Stack:    []SituationStack{{Kind: "attack_roll", DiceRoll: 3, Card: "dummy", FromDeck: true, Slot: 1}},
	}, testSet, abilitySet)
	if want := prefix + "1: an attack roll needs an attack left (R-TURN-07), and attacks is 0"; err == nil || err.Error() != want {
		t.Errorf("error %v; want %q", err, want)
	}
}

// The fields are read from a payload, and an unknown one is still an
// error.
func TestStackedDeckAttackRollJSON(t *testing.T) {
	payload := `{"sets": ["test", "abilities"],
		"setup": {"players": [{"character": "hero_a"}, {"character": "hero_b"}], "active": 0, "monsters": ["gaper"],
			"stack": [{"kind": "attack_roll", "dice_roll": 6, "card": "dummy", "owner": 0, "from_deck": true, "slot": 1}]},
		"action": {"player": 0, "kind": "pass"},
		"expect": {"allowed": true, "resolved_rolls": [6], "monsters": [{"slot": 1, "card": "gaper"}], "discards": {"monster": ["dummy"]}}}`
	if got := runExpect(t, payload); len(got) != 0 {
		t.Errorf("mismatch %q; want none: the 6 kills the revealed dummy (HP 1) and the gaper is uncovered (R-ZONE-11)", got)
	}
	bad := `{"sets": ["test"], "setup": {"players": [{"character": "hero_a"}], "active": 0, "monsters": ["gaper"],
		"stack": [{"kind": "attack_roll", "dice_roll": 6, "card": "dummy", "owner": 0, "from_deck": true, "slot": 1, "covers": "gaper"}]},
		"action": {"player": 0, "kind": "pass"}}`
	if _, err := RunSituation([]byte(bad), testSet); err == nil || !strings.Contains(err.Error(), `unknown field "covers"`) {
		t.Errorf("error %v; want an unknown field error", err)
	}
}
