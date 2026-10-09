package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed (Boss Card)
//
//	Each time this deals damage, each player loses 4¢.
var greed = engine.CardDef{
	Ref:     "greed",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 9}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this deals damage, each player loses 4¢.",
			Trigger: whenThisDealsDamage(false),
			Effects: []engine.Effect{engine.EachPlayer(engine.LoseCents(4))},
		},
	},
}
