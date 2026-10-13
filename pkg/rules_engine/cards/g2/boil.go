package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Boil (Basic Monster Card)
//
//	Each time this deals combat damage, it heals 1{HP}.
var boil = engine.CardDef{
	Ref:     "boil",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this deals combat damage, it heals 1 HP.",
		Trigger: whenThisDealsDamage(true),
		Effects: []engine.Effect{healThis(1)},
	}},
}
