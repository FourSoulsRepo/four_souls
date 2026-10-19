package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Fat Bat (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They gain +1 treasure.
var charmedFatBat = engine.CardDef{
	Ref:       "charmed_fat_bat",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        4,
	DC:        3,
	ATK:       1,
	Rewards:   []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{afterRewardsAnotherPlayer("They gain +1 treasure.", engine.GainTreasure(1))},
}
