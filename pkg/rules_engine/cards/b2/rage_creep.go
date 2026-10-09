package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Rage Creep (Basic Monster Card)
//
//	Damage this deals to the active player is also dealt to the player to their left.
var rageCreep = engine.CardDef{
	Ref:     "rage_creep",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Damage this deals to the active player is also dealt to the player to their left.",
			Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Source == self && e.Player == g.Turn.Active
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				p := c.G.PlayerAfter(c.G.Turn.Active)
				c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, c.EventAmount, engine.NoPlayer, c.Source)
			})},
		},
	},
}
