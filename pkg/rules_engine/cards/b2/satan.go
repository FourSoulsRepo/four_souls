package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Satan! (Epic Boss Card)
//
//	Each time the attacking player rolls an attack roll of 6, they choose a living player. That player dies.
var satan = engine.CardDef{
	Ref:     "satan",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    2,
	HP:      6,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 6, they choose a living player. That player dies.",
			Trigger: attackRollOnThis(6),
			Effects: []engine.Effect{killAPlayer},
		},
	},
}
