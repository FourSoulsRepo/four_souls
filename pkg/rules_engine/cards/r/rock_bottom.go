package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Rock Bottom (Passive Treasure Card)
//
//	When you gain or purchase this, put counters on it equal to the number of loot cards in your hand.
//	At the end of your turn, if you have fewer loot cards in your hand than there are counters on this, loot 2.
var rockBottom = engine.CardDef{
	Ref:    "rock_bottom",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you gain or purchase this, put counters on it equal to the number of loot cards in your hand.",
			Trigger: cardkit.WhenYouGainThis(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.AddCountersTo(c.Source, len(c.G.Players[c.Controller].Hand))
			})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, if you have fewer loot cards in your hand than there are counters on this, loot 2.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if len(c.G.Players[c.Controller].Hand) < c.G.Object(c.Source).CountersOf("") {
					c.Do(engine.Loot(2))
				}
			})},
		},
	},
}
