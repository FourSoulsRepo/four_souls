package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDistantAdmiration(t *testing.T) {
	plays := 0
	tb := setup(t, engine.SituationSetup{
		Players:   []engine.SituationPlayer{{Character: "isaac", Items: items("distant_admiration")}, {Character: "cain", Hand: items("three_cents", "bomb")}},
		LootPlays: &plays,
	})
	tb.Activate(0, "distant_admiration", 0, foe, "three_cents")
	if c := tb.G.Players[0].Cents; c != 3 {
		t.Errorf("cents = %d: the given 3 Cents! is played for free", c)
	}
	if tb.G.Players[0].ExtraLootPlays != 0 {
		t.Error("no loot play is left over")
	}
}

// A given card that cannot be played stays in hand and adds no loot
// play (19.60).
func TestDistantAdmirationUnplayable(t *testing.T) {
	plays := 1
	tb := setup(t, engine.SituationSetup{
		Players:   []engine.SituationPlayer{{Character: "isaac", Items: items("distant_admiration")}, {Character: "cain", Hand: items("butter_bean")}},
		LootPlays: &plays,
	})
	tb.Activate(0, "distant_admiration", 0, foe, "butter_bean")
	if h := tb.G.Players[0].Hand; len(h) != 1 || tb.G.Object(h[0]).Card != "butter_bean" {
		t.Errorf("hand = %v: Butter Bean! has no target, so it stays in hand", h)
	}
	if tb.G.Turn.LootPlays != 1 || tb.G.Players[0].ExtraLootPlays != 0 {
		t.Errorf("loot plays = %d + %d extra, want 1 + 0", tb.G.Turn.LootPlays, tb.G.Players[0].ExtraLootPlays)
	}
	if tb.G.FreeLoot != 0 {
		t.Error("no free play is left pending")
	}
}

// The played card uses none of the turn's loot plays (19.60).
func TestDistantAdmirationKeepsLootPlay(t *testing.T) {
	plays := 1
	tb := setup(t, engine.SituationSetup{
		Players:   []engine.SituationPlayer{{Character: "isaac", Items: items("distant_admiration")}, {Character: "cain", Hand: items("three_cents")}},
		LootPlays: &plays,
	})
	tb.Activate(0, "distant_admiration", 0, foe, "three_cents")
	if tb.G.Turn.LootPlays != 1 || tb.G.Players[0].ExtraLootPlays != 0 {
		t.Errorf("loot plays = %d + %d extra, want 1 + 0", tb.G.Turn.LootPlays, tb.G.Players[0].ExtraLootPlays)
	}
}

// A cancelled free play leaves the card in hand; played later, it uses
// a loot play (19.60).
func TestDistantAdmirationCancelled(t *testing.T) {
	plays := 1
	tb := setup(t, engine.SituationSetup{
		Players:   []engine.SituationPlayer{{Character: "isaac", Items: items("distant_admiration")}, {Character: "cain", Hand: items("bomb")}},
		LootPlays: &plays,
	})
	tb.Activate(0, "distant_admiration", 0, foe, "bomb", "cancel")
	if tb.G.FreeLoot != 0 || tb.G.Turn.LootPlays != 1 || tb.G.Players[0].ExtraLootPlays != 0 {
		t.Fatalf("free %d, loot plays %d + %d extra: the cancel leaves the plays as they were", tb.G.FreeLoot, tb.G.Turn.LootPlays, tb.G.Players[0].ExtraLootPlays)
	}
	tb.Play(0, "bomb", foe)
	if tb.G.Turn.LootPlays != 0 {
		t.Errorf("loot plays = %d: Bomb! played later uses the turn's play", tb.G.Turn.LootPlays)
	}
}
