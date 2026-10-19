package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Whipper (Basic Monster Card)
//
//	Each time a player declares an attack, the active player rolls-
//	1-2: This deals 1 damage to that player.
var whipper = engine.CardDef{
	Ref:     "whipper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player declares an attack, the active player rolls- 1-2: This deals 1 damage to that player.",
		Trigger: engine.Trigger{On: engine.EvAttackDeclared, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 2, damageActive(1)))},
	}},
}
