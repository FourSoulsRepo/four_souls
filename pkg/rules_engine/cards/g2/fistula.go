package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fistula (Boss Card)
//
//	Each time this takes damage, it gains +1{DC} till end of turn.
var fistula = engine.CardDef{
	Ref:     "fistula",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, it gains +1 DC till end of turn.",
		Trigger: whenThisTakesDamage(false),
		Effects: []engine.Effect{thisGains(engine.StatMonsterDC, 1)},
	}},
}
