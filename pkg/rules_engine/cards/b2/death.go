package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Death (Boss Card)
//
//	When this dies, the active player kills a player.
var death = engine.CardDef{
	Ref:     "death",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player kills a player.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{killAPlayer},
		},
	},
}
