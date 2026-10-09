package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestContractFromBelow(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "contract_from_below", "breakfast", "dinner"), seat("cain", "steamy_sale"),
	}}, Set)
	tb.Activate(0, "contract_from_below", 0, "breakfast", "dinner", "steamy_sale")
	if hasItem(tb.G, 0, "breakfast") || hasItem(tb.G, 0, "dinner") || !hasItem(tb.G, 0, "steamy_sale") {
		t.Error("two items for one steal did not happen")
	}
}
