package cards

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSetsStartGames(t *testing.T) {
	for _, set := range Sets() {
		for players := 2; players <= 4; players++ {
			g, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: players, Sets: []engine.CardSet{set}, BonusSouls: true})
			if err != nil {
				t.Fatalf("%s, %d players: %v", set.Code, players, err)
			}
			if err := g.CheckInvariants(); err != nil {
				t.Fatalf("%s, %d players: %v", set.Code, players, err)
			}
		}
	}
}

func TestStartingItemsExist(t *testing.T) {
	for _, set := range Sets() {
		refs := map[engine.CardRef]engine.CardDef{}
		for _, c := range set.Cards {
			refs[c.Ref] = c
		}
		for _, c := range set.Cards {
			if c.StartingItem == "" {
				continue
			}
			item, ok := refs[c.StartingItem]
			if !ok || !item.Outside {
				t.Errorf("%s: starting item %s missing or not outside the game", c.Ref, c.StartingItem)
			}
		}
	}
}

func TestFind(t *testing.T) {
	if s, ok := Find("b2"); !ok || s.Name == "" || len(s.Cards) == 0 {
		t.Fatalf("Find(b2) = %v, %v", s.Name, ok)
	}
	if _, ok := Find("nope"); ok {
		t.Fatal("Find(nope) found a set")
	}
}
