package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// War (Boss Card)
//
//	Each time this takes damage, it gains +1{ATK} till end of turn.
var war = engine.CardDef{
	Ref:     "war",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 8}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this takes damage, it gains +1 ATK till end of turn.",
			Trigger: whenThisTakesDamage(false),
			Effects: []engine.Effect{engine.GainATKThisTurn(1, engine.This)},
		},
	},
}
