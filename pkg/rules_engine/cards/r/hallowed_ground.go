package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hallowed Ground (Paid Treasure Card)
//
//	Each time a player rolls a ❷, put a counter on this.
//	{Paid Effect}Remove 3 counters from this:
//	Put the top card of the loot discard into your hand.
var hallowedGround = engine.CardDef{
	Ref:    "hallowed_ground",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a ❷, put a counter on this.",
			Trigger: engine.OnRollOf(2),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:  engine.Activated,
			Text:  "$ Remove 3 counters from this: Put the top card of the loot discard into your hand.",
			Costs: []engine.Cost{engine.RemoveCounters(3)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if pile := c.G.Discards[engine.LootDeck]; len(pile) > 0 {
					c.G.DiscardToHand(engine.LootDeck, pile[len(pile)-1], c.Controller)
				}
			})},
		},
	},
}
