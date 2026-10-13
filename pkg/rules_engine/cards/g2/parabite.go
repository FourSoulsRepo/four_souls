package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Parabite (Basic Monster Card)
//
//	Each time this takes combat damage, the attacking player has -1 to their next attack roll this turn.
var parabite = engine.CardDef{
	Ref:     "parabite",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes combat damage, the attacking player has -1 to their next attack roll this turn.",
		Trigger: whenThisTakesDamage(true),
		Effects: []engine.Effect{engine.NextAttackRoll(-1, true)},
	}},
}
