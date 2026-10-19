package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTaxForTheMighty(t *testing.T) {
	tb := roomTable(t, "tax_for_the_mighty")
	tb.G.Players[0].InPlay = append(tb.G.Players[0].InPlay, soulOf(tb, 0, "chub"))
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err == nil {
		t.Fatal("Isaac (most souls) can't pay 2¢ to attack")
	}
	tb.G.Players[0].Cents = 2
	tb.Attack("fly", 6)
	if tb.G.Players[1].Cents != 2 {
		t.Errorf("Isaac paid Cain 2¢ to attack: Cain has %d", tb.G.Players[1].Cents)
	}
}
