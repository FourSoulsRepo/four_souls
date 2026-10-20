package r

import engine "github.com/FourSoulsRepo/rules_engine"

// R Key (One-Use Treasure Card)
//
//	This enters play deactivated.
//	{Tap Effect}Destroy this. If you do, destroy all souls. Each player discards their hand and loots 3.
var rKey = engine.CardDef{
	Ref:               "r_key",
	Kind:              engine.TreasureCard,
	Copies:            1,
	Tap:               true,
	EntersDeactivated: true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Destroy this. If you do, destroy all souls. Each player discards their hand and loots 3.",
		Takes: engine.TakesOthers,
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.DestroyThis(engine.EffectFunc(func(c *engine.Ctx) {
			for _, pl := range c.G.Players {
				for _, id := range append([]engine.ObjectID(nil), pl.InPlay...) {
					if c.G.Object(id).Role == engine.RoleSoul {
						c.G.DestroyObject(pl.ID, id)
					}
				}
			}
		}), engine.EachPlayer(discardHandLoot(3)))},
	}},
}
