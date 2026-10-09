package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chub (Boss Card)
//
//	Each time the attacking player rolls an attack roll of 1, this heals 2{HP}.
var chub = engine.CardDef{
	Ref:     "chub",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 1, this heals 2 HP.",
			Trigger: attackRollOnThis(1),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.HealObject(c.Source, 2)
			})},
		},
	},
}
