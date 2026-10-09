package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lust (Boss Card)
//
//	Each time this takes combat damage, it deals 1 damage to the attacking player.
var lust = engine.CardDef{
	Ref:     "lust",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this takes combat damage, it deals 1 damage to the attacking player.",
			Trigger: whenThisTakesDamage(true),
			Effects: []engine.Effect{engine.DealDamage(1, engine.You)},
		},
	},
}
