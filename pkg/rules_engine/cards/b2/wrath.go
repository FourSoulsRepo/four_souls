package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Wrath (Boss Card)
//
//	When this dies, the active player rolls-
//	1-3: Each player takes 1 damage.
//	4-6: Each player takes 2 damage.
var wrath = engine.CardDef{
	Ref:     "wrath",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player rolls- 1-3: Each player takes 1 damage. 4-6: Each player takes 2 damage.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 3, engine.EachPlayer(engine.DealDamage(1, engine.You))).Results(4, 6, engine.EachPlayer(engine.DealDamage(2, engine.You))))},
		},
	},
}
