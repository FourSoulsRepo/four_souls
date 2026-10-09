package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Keeper Head (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❹, they gain 2¢.
var holyKeeperHead = engine.CardDef{
	Ref:     "holy_keeper_head",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 4, they gain 2¢.",
			Trigger: playerRollOf(4),
			Effects: []engine.Effect{forRoller(engine.GainCents(2))},
		},
	},
}
