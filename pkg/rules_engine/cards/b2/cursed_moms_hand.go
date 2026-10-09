package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Mom’s Hand (Cursed Monster Card)
//
//	{Curse Effect}When the active player rolls a 6, cancel everything that hasn't resolved and end the turn.
var cursedMomsHand = engine.CardDef{
	Ref:     "cursed_moms_hand",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "When the active player rolls a 6, cancel everything that hasn't resolved and end the turn.",
			Trigger: engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
				return e.Amount == 6 && e.Player == g.Turn.Active
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.EndTurnNow()
			})},
		},
	},
}
