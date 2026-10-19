package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bag Of Crafting (Eternal Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Put a counter on this.
//	{Paid Effect}Remove 4 counters from this:
//	Gain +1 treasure.
//	-Eternal- This can't be destroyed or put into discard.
var bagOfCrafting = engine.CardDef{
	Ref:     "bag_of_crafting",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "$ Discard a loot card: Put a counter on this.",
			Costs:   []engine.Cost{engine.DiscardChosen(0)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Remove 4 counters from this: Gain +1 treasure.",
			Costs:   []engine.Cost{engine.RemoveCounters(4)},
			Effects: []engine.Effect{engine.GainTreasure(1)},
		},
	},
}
