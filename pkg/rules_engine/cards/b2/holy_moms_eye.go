package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Mom’s Eye (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❷, they may recharge an item.
var holyMomsEye = engine.CardDef{
	Ref:     "holy_moms_eye",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 2, they may recharge an item.",
			Trigger: playerRollOf(2),
			Effects: []engine.Effect{forRoller(rechargeAnItem)},
		},
	},
}
