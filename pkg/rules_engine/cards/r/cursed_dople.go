package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Cursed Dople (Cursed Monster Card)
//
//	Each time this takes damage, it deals 1 damage to each non-active player.
var cursedDople = engine.CardDef{
	Ref:       "cursed_dople",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        3,
	DC:        4,
	ATK:       2,
	Rewards:   []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, it deals 1 damage to each non-active player.",
		Trigger: cardkit.WhenThisTakesDamage(false),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			for _, p := range othersOfActive(c.G) {
				c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
			}
		})},
	}},
}
