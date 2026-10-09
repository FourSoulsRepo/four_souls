package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pestilence (Boss Card)
//
//	When this dies, the active player deals 2 damage divided as they choose to any number of monsters or players.
var pestilence = engine.CardDef{
	Ref:     "pestilence",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player deals 2 damage divided as they choose to any number of monsters or players.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{damageAnything(1), damageAnything(1)},
		},
	},
}
