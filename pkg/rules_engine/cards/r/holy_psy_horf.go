package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Psy Horf (Holy/Charmed Monster Card)
//
//	At the end of each turn, the active player loots 1.
var holyPsyHorf = engine.CardDef{
	Ref:     "holy_psy_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of each turn, the active player loots 1.",
		Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }},
		Effects: []engine.Effect{engine.Loot(1)},
	}},
}
