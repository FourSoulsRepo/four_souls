package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Holy Brain (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❸, each monster gains -1{DC} till end of turn.
var holyBrain = engine.CardDef{
	Ref:     "holy_brain",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player rolls a ❸, each monster gains -1 DC till end of turn.",
		Trigger: engine.OnRollOf(3),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			for _, id := range cardkit.MonstersInPlay(c.G) {
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatMonsterDC, Player: engine.NoPlayer, Object: id, Amount: -1})
			}
		})},
	}},
}
