package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Dip (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❶, they gain 1¢.
var holyDip = engine.CardDef{
	Ref:     "holy_dip",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 1, they gain 1¢.",
			Trigger: playerRollOf(1),
			Effects: []engine.Effect{forRoller(engine.GainCents(1))},
		},
	},
}
