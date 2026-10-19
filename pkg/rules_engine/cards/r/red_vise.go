package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Vise (Room Card)
//
//	At the end of the turn, if the active player has 1 or fewer loot cards in their hand, they take 1 damage.
var redVise = engine.CardDef{
	Ref:    "red_vise",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of the turn, if the active player has 1 or fewer loot cards in their hand, they take 1 damage.",
		Trigger: atEndOfTheTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if len(c.G.Players[c.Controller].Hand) <= 1 {
				c.G.DealDamageTo(engine.Target{Player: c.Controller, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
			}
		})},
	}},
}
