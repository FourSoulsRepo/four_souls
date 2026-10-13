package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Brain (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 6, they take 1 damage.
var brain = engine.CardDef{
	Ref:     "brain",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time the attacking player rolls an attack roll of 6, they take 1 damage.",
		Trigger: attackRollOnThis(6),
		Effects: []engine.Effect{damageActive(1)},
	}},
}
