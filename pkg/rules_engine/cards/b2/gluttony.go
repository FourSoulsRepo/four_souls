package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gluttony (Boss Card)
//
//	Each time this takes combat damage on an attack roll of 6, deal 1 damage to the player to the active player's left.
var gluttony = engine.CardDef{
	Ref:     "gluttony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time this takes combat damage on an attack roll of 6, deal 1 damage to the player to the active player's left.",
			Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Object == self && e.Text == "combat" && g.Attack.LastRoll == 6
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				p := c.G.PlayerAfter(c.G.Turn.Active)
				c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
			})},
		},
	},
}
