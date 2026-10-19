package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Splash Damage (Room Card)
//
//	Each time a player deals damage to a monster, they deal 1 damage to the player to their left.
var splashDamage = engine.CardDef{
	Ref:       "splash_damage",
	Kind:      engine.RoomCard,
	Copies:    1,
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time a player deals damage to a monster, they deal 1 damage to the player to their left.",
		Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Player == engine.NoPlayer && e.By > 0 && e.Object != 0 && g.Object(e.Object).Role == engine.RoleMonster
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			by := engine.PlayerID(c.EventBy - 1)
			c.G.DealDamageTo(engine.Target{Player: c.G.PlayerAfter(by), IsPlayer: true}, 1, by, c.Source)
		})},
	}},
}
