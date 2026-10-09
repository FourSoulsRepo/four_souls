package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Bloat (Boss Card)
//
//	Each time this deals combat damage, it deals 1 damage to each non-active player.
var theBloat = engine.CardDef{
	Ref:     "the_bloat",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this deals combat damage, it deals 1 damage to each non-active player.",
			Trigger: whenThisDealsDamage(true),
			Effects: []engine.Effect{engine.EachOtherPlayer(engine.DealDamage(1, engine.You))},
		},
	},
}
