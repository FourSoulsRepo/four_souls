package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dank Globin (Basic Monster Card)
//
//	When this dies, the active player forces a player to discard 2 loot cards.
var dankGlobin = engine.CardDef{
	Ref:     "dank_globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player forces a player to discard 2 loot cards.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(forcedDiscard(2), forcedDiscardQuestions(2)...)},
		},
	},
}
