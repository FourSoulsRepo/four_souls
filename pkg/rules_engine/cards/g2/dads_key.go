package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dad’s Key (Paid Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Put the top card of the monster deck in a monster slot not being attacked.
var dadsKey = engine.CardDef{
	Ref:    "dads_key",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "Discard a loot card: Put the top card of the monster deck in a monster slot not being attacked.",
		Costs:   []engine.Cost{engine.DiscardChosen(0)},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
		Effects: []engine.Effect{coverASlot(false)},
	}},
}
