package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// MOTHER! (Epic Boss Card)
//
//	Each time this takes damage, it gains +1{DC} till end of turn.
//	When this dies, expand monster slots by 2.
var mother = engine.CardDef{
	Ref:     "mother",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      8,
	DC:      1,
	ATK:     8,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this takes damage, it gains +1 DC till end of turn.",
			Trigger: cardkit.WhenThisTakesDamage(false),
			Effects: []engine.Effect{thisGains(engine.StatMonsterDC, 1)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, expand monster slots by 2.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandMonsters(2) })},
		},
	},
}
