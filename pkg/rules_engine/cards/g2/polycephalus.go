package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Polycephalus (Boss Card)
//
//	Each time this takes damage, the next attack roll this turn has -2.
var polycephalus = engine.CardDef{
	Ref:     "polycephalus",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, the next attack roll this turn has -2.",
		Trigger: whenThisTakesDamage(false),
		Effects: []engine.Effect{engine.NextAttackRoll(-2, false)},
	}},
}
