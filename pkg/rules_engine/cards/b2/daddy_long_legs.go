package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Daddy Long Legs (Boss Card)
//
//	Each time the attacking player rolls an attack roll of 1, each monster gains +1{DC} till end of turn.
var daddyLongLegs = engine.CardDef{
	Ref:     "daddy_long_legs",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 1, each monster gains +1 DC till end of turn.",
			Trigger: attackRollOnThis(1),
			Effects: []engine.Effect{eachMonsterGains(engine.StatMonsterDC)},
		},
	},
}
