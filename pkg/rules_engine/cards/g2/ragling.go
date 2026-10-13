package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ragling (Basic Monster Card)
//
//	When this dies, after gaining rewards, the active player rolls-
//	1 or 6: Put this on top of the Monster Deck.
var ragling = engine.CardDef{
	Ref:     "ragling",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, after gaining rewards, the active player rolls- 1 or 6: Put this on top of the monster deck.",
		Trigger: engine.AfterThisRewards(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 1, backOnTop).Results(6, 6, backOnTop))},
	}},
}
