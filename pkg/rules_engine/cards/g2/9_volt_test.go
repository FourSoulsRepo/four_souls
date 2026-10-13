package g2

import (
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
