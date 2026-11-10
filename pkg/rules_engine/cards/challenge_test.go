package cards

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
)

// ADR 028: the challenge set is no playable set; its cards are never dealt.
func TestChallengeSetIsNotDealt(t *testing.T) {
	if _, ok := Find(challenge.Set.Code); ok {
		t.Fatal("the challenge set is listed as a playable set")
	}
	sets := append(Sets(), challenge.Set)
	var refs []engine.CardRef
	for _, d := range challenge.Set.Cards {
		refs = append(refs, d.Ref)
	}
	for seed := range uint64(5) {
		g, _, err := engine.NewGame(engine.Setup{Seed: seed, Players: 2, Sets: sets})
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range g.Objects {
			if slices.Contains(refs, o.Card) {
				t.Fatalf("seed %d: %s was dealt, zone %+v", seed, o.Card, o.Zone)
			}
		}
	}
}

// R-KRAMPUS-03: Pilfering Krampus by difficulty; the rule cards have
// their own kind (ADR 028).
func TestKrampusCards(t *testing.T) {
	defs := map[engine.CardRef]engine.CardDef{}
	for _, d := range challenge.Set.Cards {
		defs[d.Ref] = d
	}
	for ref, hp := range map[engine.CardRef]int{
		"pilfering_krampus_norm":  5,
		"pilfering_krampus_hard":  6,
		"pilfering_krampus_ultra": 6,
		"pilfering_krampus_comp":  5, // the competitive boss (14.28)
	} {
		d, ok := defs[ref]
		want := []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}}
		if !ok || d.Kind != engine.MonsterCard || !d.Outside || d.HP != hp || d.DC != 4 || d.ATK != 2 ||
			d.Soul != 1 || !slices.Equal(d.Rewards, want) {
			t.Errorf("%s = %+v", ref, d)
		}
	}
	for _, diff := range []string{"norm", "hard", "ultra", "comp"} {
		ref := engine.CardRef("how_the_krampus_stole_christmas_" + diff)
		if d, ok := defs[ref]; !ok || d.Kind != engine.ChallengeCard {
			t.Errorf("%s = %+v", ref, d)
		}
	}
	if _, deck := engine.ChallengeCard.Deck(); deck {
		t.Error("a rule card has a deck")
	}
	if len(defs) != 8 {
		t.Errorf("%d challenge cards, want 8", len(defs))
	}
}
