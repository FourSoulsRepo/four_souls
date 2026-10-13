package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Knight (Basic Monster Card)
//
//	The attacking player also deals combat damage to this on attack rolls of 1.
var knight = engine.CardDef{
	Ref:     "knight",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      6,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "The attacking player also deals combat damage to this on attack rolls of 1.",
		Trigger: attackRollOnThis(1),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.DealDamageTo(engine.Target{Object: c.Source}, c.G.PlayerATK(c.G.Turn.Active), c.G.Turn.Active, c.G.Players[c.G.Turn.Active].Character)
		})},
	}},
}
