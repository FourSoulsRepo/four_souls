package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Dinga (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❻, they heal 1{HP}.
var holyDinga = engine.CardDef{
	Ref:     "holy_dinga",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 6, they heal 1 HP.",
			Trigger: playerRollOf(6),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.HealPlayer(c.EventPlayer, 1) })},
		},
	},
}
