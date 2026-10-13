package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Krampus (Boss Card)
//
//	The active player must attack this once each turn if able.
var krampus = engine.CardDef{
	Ref:     "krampus",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "The active player must attack this once each turn if able.",
			Trigger: engine.Trigger{On: engine.EvStartOfTurn, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }},
			Effects: []engine.Effect{mustAttackThis},
		},
		{
			Kind:    engine.Triggered,
			Text:    "The active player must attack this once each turn if able.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{mustAttackThis},
		},
	},
}
