package rulesengine

import (
	"encoding/json"
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
