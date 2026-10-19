package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Keeper’s Penny (Active Treasure Card)
//
//	When you gain or purchase this, gain 9¢.
//	{Tap Effect}Put a counter on this. Gain ¢ equal to the number of counters on this, then give this card to another player.
var keepersPenny = engine.CardDef{
	Ref:    "keepers_penny",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you gain or purchase this, gain 9¢.",
			Trigger: cardkit.WhenYouGainThis(),
			Effects: []engine.Effect{engine.GainCents(9)},
		},
		{
			Kind:  engine.Activated,
			Text:  "↷: Put a counter on this. Gain ¢ equal to the number of counters on this, then give this card to another player.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.AddCounters(1), engine.EffectFunc(func(c *engine.Ctx) {
				c.Do(engine.GainCents(c.G.Object(c.Source).CountersOf("")))
			}), cardkit.GiveThisAway},
		},
	},
}
