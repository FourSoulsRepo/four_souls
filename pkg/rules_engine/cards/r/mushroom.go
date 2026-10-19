package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Mushroom (Basic Monster Card)
//
//	Each time this deals combat damage to the attacking player, it deals 1 damage to each other player.
var mushroom = engine.CardDef{
	Ref:     "mushroom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this deals combat damage to the attacking player, it deals 1 damage to each other player.",
		Trigger: cardkit.WhenThisDealsDamage(true),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			for _, p := range othersOfActive(c.G) {
				c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
			}
		})},
	}},
}
