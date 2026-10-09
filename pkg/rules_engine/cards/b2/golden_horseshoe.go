package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Golden Horseshoe (Trinket Card)
//
//	At the start of your turn, look at the top card of the treasure deck, you may put it on the bottom.
//	-Trinket- This loot becomes an item under your control when it resolves.
var goldenHorseshoe = engine.CardDef{
	Ref:     "golden_horseshoe",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, look at the top card of the treasure deck. You may put it on the bottom.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{lookMayBottom(engine.TreasureDeck)},
		},
	},
}
