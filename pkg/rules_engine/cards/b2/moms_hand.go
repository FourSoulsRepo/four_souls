package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Hand (Basic Monster Card)
//
//	When the attacking player rolls an attack roll of 6, cancel everything that hasn't resolved and end the turn.
var momsHand = engine.CardDef{
	Ref:     "moms_hand",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When the attacking player rolls an attack roll of 6, cancel everything that hasn't resolved and end the turn.",
			Trigger: attackRollOnThis(6),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.EndTurnNow()
			})},
		},
	},
}
