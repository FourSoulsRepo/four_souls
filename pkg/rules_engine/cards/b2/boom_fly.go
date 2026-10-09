package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Boom Fly (Basic Monster Card)
//
//	When this dies, it deals 1 damage to each player.
var boomFly = engine.CardDef{
	Ref:     "boom_fly",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, it deals 1 damage to each player.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EachPlayer(engine.DealDamage(1, engine.You))},
		},
	},
}
