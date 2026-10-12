package rulesengine

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// ids collects every object ID mentioned in a view.
func viewIDs(v GameView) map[ObjectID]bool {
	ids := map[ObjectID]bool{}
	add := func(c CardView) { ids[c.ID] = true }
	for _, p := range v.Players {
		add(p.Character)
		for _, c := range p.InPlay {
			add(c)
		}
		for _, c := range p.Hand {
			add(c)
		}
	}
	for _, s := range append(v.Shop, v.Monsters...) {
		if s.Top != nil {
			add(*s.Top)
		}
		for _, c := range s.Covered {
			add(c)
		}
	}
	for _, d := range v.Discards {
		for _, c := range d {
			add(c)
		}
	}
	return ids
}

func TestViewsHideWhatTheyMust(t *testing.T) {
	g, _, err := NewGame(Setup{Seed: 33, Players: 3, Sets: []CardSet{testSet, abilitySet}})
	if err != nil {
		t.Fatal(err)
	}
	for step := range 300 {
		if g.Over {
			break
		}
		for _, pl := range g.Players {
			pv := g.View(Viewer{Kind: ViewPlayer, Player: pl.ID})
			ids := viewIDs(pv)
			for _, other := range g.Players {
				for _, h := range other.Hand {
					if ids[h] != (other.ID == pl.ID) {
						t.Fatalf("step %d: player %d sees hand card %d of player %d: %v (A-07)", step, pl.ID, h, other.ID, ids[h])
					}
				}
				if pv.Players[other.ID].HandSize != len(other.Hand) {
					t.Fatal("hand sizes are public (R-ZONE-08)")
				}
			}
			for d := range deckCount {
				for _, id := range g.Decks[d] {
					if ids[id] {
						t.Fatalf("a deck card is visible (R-ZONE-04)")
					}
				}
			}
		}
		judge := viewIDs(g.View(Viewer{Kind: ViewJudge}))
		spectator := viewIDs(g.View(Viewer{Kind: ViewSpectator}))
		for _, pl := range g.Players {
			for _, h := range pl.Hand {
				if !judge[h] {
					t.Fatal("the judge sees every hand (SP-02)")
				}
				if spectator[h] {
					t.Fatal("spectators never see hands (SP-01)")
				}
			}
		}
		data, err := json.Marshal(g.View(Viewer{Kind: ViewJudge}))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), `"rng"`) || strings.Contains(string(data), `"inc"`) {
			t.Fatal("the RNG must never be in a view")
		}
		allowed := g.Allowed(g.Prompt().Player)
		in := allowed[(step*5)%len(allowed)]
		if in.Kind == IntentDiscard {
			in.Objects = in.Objects[:g.Prompt().Count]
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatal(err)
		}
	}
}

func TestChoiceOptionsAreHiddenFromOthers(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	g.ask(Choice{Purpose: ChoosePenaltyLoot, Player: a, Rule: "R-DEATH-14", Objects: g.Players[a].Hand}, g.labels(g.Players[a].Hand))
	if len(g.View(Viewer{Kind: ViewPlayer, Player: a}).Waiting.Options) == 0 {
		t.Error("the asked player sees the options")
	}
	if len(g.View(Viewer{Kind: ViewPlayer, Player: g.next(a)}).Waiting.Options) != 0 {
		t.Error("other players must not see options that may name hidden cards")
	}
	if len(g.View(Viewer{Kind: ViewJudge}).Waiting.Options) == 0 {
		t.Error("the judge sees the options")
	}
}

func TestPrivateEventsAreFiltered(t *testing.T) {
	events := []Event{{Kind: EvLooted, Player: 1, Object: 42, Card: "penny", Private: true}, {Kind: EvGainedCents, Player: 1, Amount: 3}}
	other := FilterEvents(events, Viewer{Kind: ViewPlayer, Player: 0})
	if other[0].Card != "" || other[0].Object != 0 || other[1].Amount != 3 {
		t.Errorf("other player sees %+v", other)
	}
	own := FilterEvents(events, Viewer{Kind: ViewPlayer, Player: 1})
	if own[0].Card != "penny" {
		t.Error("the owner sees their looted card")
	}
	if FilterEvents(events, Viewer{Kind: ViewJudge})[0].Card != "penny" {
		t.Error("the judge sees looted cards")
	}
}

func TestChoiceOptionsPointAtTheTable(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	m := putMonster(g, 0, "gaper")
	attack(t, g)
	v := g.View(Viewer{Kind: ViewPlayer, Player: a})
	if len(v.Options) != len(v.Waiting.Options) {
		t.Fatalf("%d option views for %d options", len(v.Options), len(v.Waiting.Options))
	}
	if v.Options[0].Object != m {
		t.Errorf("option 0 points at %+v, want the gaper %d", v.Options[0], m)
	}
	if last := v.Options[len(v.Options)-1]; !last.Deck {
		t.Errorf("the last option %+v should be the monster deck", last)
	}
	if other := g.View(Viewer{Kind: ViewPlayer, Player: g.next(a)}); other.Options != nil {
		t.Error("others must not see what the options point at")
	}
}

func TestViewCountsWhatIsLeft(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	id := g.Players[a].Character
	g.Object(id).addCounters("charge", 2)
	v := g.View(Viewer{Kind: ViewPlayer, Player: g.next(a)})
	me := v.Players[a]
	if me.LootPlays != g.lootPlaysFor(a) || me.Attacks != g.attacksLeft() || me.Purchases != g.purchasesLeft() {
		t.Errorf("active player left %d/%d/%d, want %d/%d/%d",
			me.LootPlays, me.Attacks, me.Purchases, g.lootPlaysFor(a), g.attacksLeft(), g.purchasesLeft())
	}
	if other := v.Players[g.next(a)]; other.Attacks != 0 || other.Purchases != 0 {
		t.Errorf("a waiting player has attacks %d and purchases %d, want none", other.Attacks, other.Purchases)
	}
	if c := me.Character.Counters; len(c) != 1 || c[0].Count != 2 {
		t.Errorf("counters %+v, want 2 charge", c)
	}
}

func TestViewNamesWhoPassedAndWhoMayRespond(t *testing.T) {
	g := toAction(t, 4)
	a := g.Turn.Active
	b := g.next(a)
	c, d := g.next(b), g.next(g.next(b))
	g.push(StackItem{Kind: StackRoll, Controller: b, Label: "x"})
	pass(t, g)
	pass(t, g)
	// Every viewer sees the same lists (V-11).
	for _, v := range []Viewer{{Kind: ViewPlayer, Player: a}, {Kind: ViewPlayer, Player: c}, {Kind: ViewSpectator}} {
		gv := g.View(v)
		if !slices.Equal(gv.Passed, []PlayerID{b, c}) || !slices.Equal(gv.Responding, []PlayerID{d, a}) {
			t.Errorf("viewer %+v: passed %v, responding %v; want [%d %d] and [%d %d]",
				v, gv.Passed, gv.Responding, b, c, d, a)
		}
	}
	g.Waiting = Prompt{Kind: PromptChoose, Player: a}
	if gv := g.View(Viewer{Kind: ViewSpectator}); gv.Passed != nil || gv.Responding != nil {
		t.Errorf("a choice is not a priority window: passed %v, responding %v", gv.Passed, gv.Responding)
	}
}

func TestGivenCardIsSeenByBothHands(t *testing.T) {
	g := toAction(t, 3)
	from, to, other := PlayerID(0), PlayerID(1), PlayerID(2)
	id := g.Players[from].Hand[0]
	card := g.Object(id).Card
	g.events = nil
	g.GiveHandCard(from, to, id)
	ev := g.takeEvents()
	for _, c := range []struct {
		v    Viewer
		sees bool
	}{
		{Viewer{Kind: ViewPlayer, Player: from}, true},
		{Viewer{Kind: ViewPlayer, Player: to}, true},
		{Viewer{Kind: ViewPlayer, Player: other}, false},
		{Viewer{Kind: ViewSpectator}, false},
		{Viewer{Kind: ViewJudge}, true},
	} {
		got := FilterEvents(ev, c.v)
		if len(got) != 1 || got[0].Kind != EvGaveCard {
			t.Fatalf("viewer %+v: events %+v, want one gave_card", c.v, got)
		}
		if seen := got[0].Card == card; seen != c.sees {
			t.Errorf("viewer %+v sees the given card: %v, want %v (V-14)", c.v, seen, c.sees)
		}
		if got[0].Text != "0" {
			t.Errorf("viewer %+v: giver %q, want everyone to know it", c.v, got[0].Text)
		}
	}
}
