package fsp2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Headless Horseman (Boss Card)
//
//	The first time this would die each turn, prevent death. This heals 2{HP} and gains +1{DC} and -1{ATK} till end of turn.
var headlessHorseman = engine.CardDef{
	Ref:     "headless_horseman",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "The first time this would die each turn, prevent death. This heals 2 HP and gains +1 DC and -1 ATK till end of turn.",
		Trigger: engine.Trigger{On: engine.EvDeathPending, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Player == engine.NoPlayer && e.Object == self && !cardkit.Used(g, self)
		}},
		Effects: []engine.Effect{
			engine.EffectFunc(func(c *engine.Ctx) { cardkit.MarkUsed(c.G, c.Source) }), engine.PreventThisDeath(), healThis(2),
			thisGains(engine.StatMonsterDC, 1), thisGains(engine.StatMonsterATK, -1),
		},
	}},
}
