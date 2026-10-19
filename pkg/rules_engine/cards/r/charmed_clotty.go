package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Clotty (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They loot 1.
var charmedClotty = engine.CardDef{
	Ref:       "charmed_clotty",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        2,
	DC:        3,
	ATK:       1,
	Rewards:   []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{afterRewardsAnotherPlayer("They loot 1.", engine.Loot(1))},
}
