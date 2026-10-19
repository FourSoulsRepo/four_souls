package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mama Mega (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, deal 3 damage to each monster and player.
var mamaMega = engine.CardDef{
	Ref:    "mama_mega",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Destroy this. If you do, deal 3 damage to each monster and player.",
		Costs: []engine.Cost{engine.Tap()},
		// Players' damage goes on the stack first; the monsters' goes on it
		// last, so monsters take theirs first (R-MECH-30).
		Effects: []engine.Effect{engine.DestroyThis(engine.EffectFunc(func(c *engine.Ctx) {
			order := c.G.TurnOrderFrom(c.Controller)
			for i := len(order) - 1; i >= 0; i-- {
				c.G.DealDamageTo(engine.Target{Player: order[i], IsPlayer: true}, 3, c.Controller, c.Source)
			}
		}), engine.EachMonsterTakesDamage(3))},
	}},
}
