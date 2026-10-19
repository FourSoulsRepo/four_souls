package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed Looms (Room Card)
//
//	At the end of the turn, if the active player didn't purchase a shop item, they discard their hand.
var greedLooms = engine.CardDef{
	Ref:    "greed_looms",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of the turn, if the active player didn't purchase a shop item, they discard their hand.",
		Trigger: atEndOfTheTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if c.G.Turn.ShopBought == 0 {
				c.Do(discardHandLoot(0))
			}
		})},
	}},
}
