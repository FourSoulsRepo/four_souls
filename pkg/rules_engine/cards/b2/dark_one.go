package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark One (Boss Card)
//
//	Each time this takes damage, it gains +1{ATK} till end of turn.
var darkOne = engine.CardDef{
	Ref:     "dark_one",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this takes damage, it gains +1 ATK till end of turn.",
			Trigger: whenThisTakesDamage(false),
			Effects: []engine.Effect{engine.GainATKThisTurn(1, engine.This)},
		},
	},
}
