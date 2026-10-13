package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestQuestionmarkCard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"questionmark_card"}},
		{Character: "cain", Items: []engine.CardRef{"mystery_sack"}, Deactivated: []engine.CardRef{"mystery_sack"}},
	}})
	tb.G.ForceRolls(3)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "questionmark_card")}}); err != nil {
		t.Fatal(err)
	}
	opts := tb.Options()
	if len(opts) != 2 {
		t.Fatalf("options = %v, want the sack's ↷ and cancel", opts)
	}
	tb.Choose(opts[0])
	tb.Settle()
	if c := tb.G.Players[0].Cents; c != 4 {
		t.Errorf("cents = %d: the copied sack ability should give Isaac 4¢ on a 3", c)
	}
	if tb.G.Object(tb.Find(1, "mystery_sack")).Charged {
		t.Error("copying does not recharge or tap the sack")
	}
	if !inDiscard(tb.G, engine.LootDeck, "questionmark_card") {
		t.Error("the ? Card goes to the loot discard")
	}
}

func TestQuestionmarkCardNeedsAnItem(t *testing.T) {
	tb := lootTable(t, "questionmark_card")
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "questionmark_card")}}); err == nil {
		t.Error("without a ↷ item in play it can't be played")
	}
}
