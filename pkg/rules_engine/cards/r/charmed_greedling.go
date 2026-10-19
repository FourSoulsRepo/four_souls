package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Greedling (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They gain 10¢.
var charmedGreedling = engine.CardDef{
	Ref:       "charmed_greedling",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        1,
	DC:        6,
	ATK:       0,
	Rewards:   []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	Abilities: []engine.Ability{afterRewardsAnotherPlayer("They gain 10¢.", engine.GainCents(10))},
}
