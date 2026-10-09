package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Psy Horf (Basic Monster Card)
//
//	When this dies, the active player recharges each item they control.
var psyHorf = engine.CardDef{
	Ref:     "psy_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player recharges each item they control.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.RechargeItemsOf(engine.You)},
		},
	},
}
