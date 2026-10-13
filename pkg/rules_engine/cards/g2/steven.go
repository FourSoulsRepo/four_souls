package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Steven (Boss Card)
//
//	Each time the attacking player rolls an attack roll of 2, they take 2 damage.
//	Each time the attacking player rolls an attack roll of 6, cancel everything that hasn't resolved and end the turn.
var steven = engine.CardDef{
	Ref:     "steven",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 2, they take 2 damage.",
			Trigger: attackRollOnThis(2),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.DealDamageTo(engine.Target{Player: c.G.Turn.Active, IsPlayer: true}, 2, engine.NoPlayer, c.Source)
			})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 6, cancel everything that hasn't resolved and end the turn.",
			Trigger: attackRollOnThis(6),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.EndTurnNow() })}, // R-MECH-52
		},
	},
}
