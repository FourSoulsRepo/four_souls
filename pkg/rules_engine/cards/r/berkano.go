package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Berkano (Pill/Rune Card)
//
//	Destroy an item. If that item was controlled by a player, they steal an item from the shop.
var berkano = engine.CardDef{
	Ref:    "berkano",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Destroy an item. If that item was controlled by a player, they steal an item from the shop.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetItem)},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			id := c.Targets[0].Object
			p := c.G.Object(id).Controller
			// Only a destroyed item gives the steal (Q-19, owner).
			if c.G.DestroyObject(p, id) && p != engine.NoPlayer {
				c.Do(stealFromShop(p, false))
			}
		})},
	}},
}
