package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Scolex (Boss Card)
//
//	Each time this deals combat damage to a player, they discard a loot card.
var scolex = engine.CardDef{
	Ref:     "scolex",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this deals combat damage to a player, they discard a loot card.",
			Trigger: whenThisDealsDamage(true),
			Effects: []engine.Effect{rollerDiscards},
		},
	},
}
