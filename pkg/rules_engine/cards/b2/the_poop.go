package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Poop (Paid Treasure Card)
//
//	Each time you take damage, put a counter on this.
//	{Paid Effect}Remove a counter from this:
//	Prevent the next 1 damage you would take this turn.
var thePoop = engine.CardDef{
	Ref:    "the_poop",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, put a counter on this.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "Remove a counter from this: Prevent the next 1 damage you would take this turn.",
			Costs:   []engine.Cost{engine.RemoveCounters(1)},
			Effects: []engine.Effect{engine.PreventDamage(1, engine.You)},
		},
	},
}
