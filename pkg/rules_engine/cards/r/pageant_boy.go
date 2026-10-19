package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pageant Boy (One-Use Treasure Card)
//
//	{Paid Effect}Destroy this: Roll-
//	Gain x¢, where x is 4 times the result.
var pageantBoy = engine.CardDef{
	Ref:    "pageant_boy",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "$ Destroy this: Roll- Gain x¢, where x is 4 times the result.",
		Costs:   []engine.Cost{engine.DestroySelf()},
		Effects: []engine.Effect{engine.Roll(pageantTable())},
	}},
}

func pageantTable() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.GainCents(4*r))
	}
	return t
}
