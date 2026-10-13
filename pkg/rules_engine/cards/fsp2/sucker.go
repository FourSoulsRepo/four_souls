package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sucker (Basic Monster Card)
//
//	When this dies, it deals 1 damage to each player.
var sucker = engine.CardDef{
	Ref:     "sucker",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, it deals 1 damage to each player.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{eachPlayerTakes(1)},
	}},
}
