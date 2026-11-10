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

// ADR 028, Decision 3: the challenge list; each level names a boss and
// a rule card of the challenge's own set.
func TestChallenges(t *testing.T) {
	c, ok := FindChallenge("krampus")
	if !ok || c != Challenges()[0] {
		t.Fatal("no krampus challenge")
	}
	if _, ok := FindChallenge("nope"); ok {
		t.Error("found an unknown challenge")
	}
	sets := map[string]engine.CardSet{challenge.Set.Code: challenge.Set}
	for _, c := range Challenges() {
		set, ok := sets[c.Set]
		if !ok || c.ID == "" || len(c.Levels) == 0 || c.Setup == nil {
			t.Errorf("%s: set %q, %d levels", c.ID, c.Set, len(c.Levels))
		}
		kinds := map[engine.CardRef]engine.CardKind{}
		for _, d := range set.Cards {
			kinds[d.Ref] = d.Kind
		}
		var diffs []string
		for _, l := range c.Levels {
			diffs = append(diffs, l.Difficulty)
			if k, ok := kinds[l.Boss]; !ok || k != engine.MonsterCard {
				t.Errorf("%s %s: boss %q", c.ID, l.Difficulty, l.Boss)
			}
			if k, ok := kinds[l.Card]; !ok || k != engine.ChallengeCard {
				t.Errorf("%s %s: rule card %q", c.ID, l.Difficulty, l.Card)
			}
		}
		if c.ID == "krampus" && !slices.Equal(diffs, []string{engine.DifficultyNormal, engine.DifficultyHard, engine.DifficultyUltra}) {
			t.Errorf("krampus difficulties %v (R-KRAMPUS-02)", diffs)
		}
	}
}

// R-KRAMPUS-04, R-KRAMPUS-05: a Krampus game starts at every
// difficulty, with its boss in a final boss slot and the 4 presents.
func TestKrampusPlayable(t *testing.T) {
	for _, d := range krampusLevels {
		g, _, err := engine.NewGame(krampusSetup(1, d))
		if err != nil {
			t.Fatalf("%s: %v", d, err)
		}
		if len(g.Bosses) != 1 || len(g.Challenge.Presents) != 4 {
			t.Errorf("%s: %d bosses, %d presents", d, len(g.Bosses), len(g.Challenge.Presents))
		}
		checkInvariants(t, g)
	}
}
