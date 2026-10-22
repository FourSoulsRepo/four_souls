package rulesengine

import (
	"slices"
	"testing"
)

// TestThreatens: which stack items may hurt player 1 (15.84 – 15.86).
func TestThreatens(t *testing.T) {
	set := CardSet{Name: "test", Cards: append(slices.Clone(testSet.Cards),
		CardDef{Ref: "taker", Kind: LootCard, Abilities: []Ability{{Kind: LootAbility, Takes: TakesOthers}}},
		CardDef{Ref: "bird", Kind: TreasureCard, Abilities: []Ability{{Kind: Triggered, Takes: TakesEventPlayer}}},
	)}
	g, err := BuildTable(SituationSetup{Players: []SituationPlayer{
		{Character: "hero_a"}, {Character: "hero_b", Items: []CardRef{"trinket"}, Hand: []CardRef{"penny"}}, {Character: "hero_c"},
	}}, set)
	if err != nil {
		t.Fatal(err)
	}
	item, card := ObjectID(-1), g.Players[1].Hand[0]
	for _, id := range g.Players[1].InPlay {
		if g.Objects[id].Card == "trinket" {
			item = id
		}
	}
	if item < 0 {
		t.Fatal("player 1 has no trinket")
	}
	taker := AbilityRef{Card: "taker"}
	bird := AbilityRef{Card: "bird"}
	cases := []struct {
		name string
		it   StackItem
		want bool
	}{
		{"targets them", StackItem{Controller: 0, Targets: []Chosen{{Kind: TargetPlayer, Player: 1}}}, true},
		{"targets their item", StackItem{Controller: 0, Targets: []Chosen{{Kind: TargetItem, Object: item}}}, true},
		{"targets their hand card", StackItem{Controller: 2, Targets: []Chosen{{Kind: TargetYourHandCard, Object: card}}}, true},
		{"damage to them", StackItem{Kind: StackDamage, Controller: 0, Target: Target{Player: 1, IsPlayer: true}}, true},
		{"damage to their item", StackItem{Kind: StackDamage, Controller: 0, Target: Target{Object: item}}, true},
		{"takes from others", StackItem{Controller: 0, Ability: taker}, true},
		{"takes from the roller: them", StackItem{Controller: 2, Ability: bird, EventPlayer: 1}, true},
		{"takes from the roller: another", StackItem{Controller: 2, Ability: bird, EventPlayer: 0}, false},
		{"targets another player", StackItem{Controller: 0, Targets: []Chosen{{Kind: TargetPlayer, Player: 2}}}, false},
		{"their own", StackItem{Controller: 1, Ability: taker, Targets: []Chosen{{Kind: TargetPlayer, Player: 1}}}, false},
		{"nobody's", StackItem{Kind: StackDeath, Controller: NoPlayer, Target: Target{Player: 1, IsPlayer: true}}, false},
		{"a spent cost", StackItem{Controller: 0, Targets: []Chosen{{Kind: TargetItem, Object: item, Spent: true}}}, false},
		{"no ability", StackItem{Kind: StackRoll, Controller: 0}, false},
	}
	for _, c := range cases {
		if got := g.Threatens(c.it, 1); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}
