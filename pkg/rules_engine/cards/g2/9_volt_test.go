package g2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func Test9Volt(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"9_volt", "the_bible", "fruitcake", "the_butter_bean"}, Deactivated: []engine.CardRef{"the_bible", "fruitcake", "the_butter_bean"}},
		seat("cain"),
	}})
	tb.EndTurn("the_bible", "the_butter_bean")
	for card, want := range map[engine.CardRef]bool{"the_bible": true, "fruitcake": false, "the_butter_bean": true} {
		if got := tb.G.Object(tb.Find(0, card)).Charged; got != want {
			t.Errorf("%s charged = %v, want %v", card, got, want)
		}
	}
}

// Only spent active items are offered: not a charged one, and not a
// passive one even when it is spent (19.146).
func Test9VoltOffersOnlyRechargeableItems(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"9_volt", "the_bible", "the_butter_bean", "breakfast"}, Deactivated: []engine.CardRef{"the_bible", "breakfast"}},
		seat("cain"),
	}})
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	for i := 0; tb.G.Prompt().Kind != engine.PromptChoose; i++ {
		if i == 50 || tb.G.Turn.Active != 0 {
			t.Fatalf("no question at the end of the turn, prompt %+v", tb.G.Prompt())
		}
		tb.Pass(int(tb.G.Prompt().Player))
	}
	if got := tb.Options(); !slices.Contains(got, "the_bible") || slices.Contains(got, "the_butter_bean") || slices.Contains(got, "breakfast") {
		t.Errorf("options = %v, want the spent Bible only", got)
	}
}
