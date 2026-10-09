package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Swarm Of Flies (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 5, they take 1 damage.
var swarmOfFlies = engine.CardDef{
	Ref:     "swarm_of_flies",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 5, they take 1 damage.",
			Trigger: attackRollOnThis(5),
			Effects: []engine.Effect{engine.DealDamage(1, engine.You)},
		},
	},
}
