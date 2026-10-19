package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Red Ghost (Basic Monster Card)
//
//	Each time this takes damage, the active player chooses another player. That player takes 1 damage.
var redGhost = engine.CardDef{
	Ref:       "red_ghost",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        3,
	DC:        4,
	ATK:       1,
	Rewards:   []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, the active player chooses another player. That player takes 1 damage.",
		Trigger: cardkit.WhenThisTakesDamage(false),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ps := othersOfActive(c.G); a[0] >= 0 && a[0] < len(ps) {
				c.G.DealDamageTo(engine.Target{Player: ps[a[0]], IsPlayer: true}, 1, engine.NoPlayer, c.Source)
			}
		}, engine.Question{Text: "Which player takes 1 damage?", Options: func(c *engine.Ctx, _ []int) []string {
			return playerLabels(c.G, othersOfActive(c.G))
		}})},
	}},
}
