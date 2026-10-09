package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Squirt (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❺, they loot 1.
var holySquirt = engine.CardDef{
	Ref:     "holy_squirt",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 5, they loot 1.",
			Trigger: playerRollOf(5),
			Effects: []engine.Effect{forRoller(engine.Loot(1))},
		},
	},
}
