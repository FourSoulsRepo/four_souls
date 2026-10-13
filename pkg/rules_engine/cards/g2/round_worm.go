package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Round Worm (Basic Monster Card)
//
//	Each time this deals Combat Damage, the attacking player has +1 to their next attack roll this turn.
var roundWorm = engine.CardDef{
	Ref:     "round_worm",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this deals combat damage, the attacking player has +1 to their next attack roll this turn.",
		Trigger: whenThisDealsDamage(true),
		Effects: []engine.Effect{engine.NextAttackRoll(1, true)},
	}},
}
