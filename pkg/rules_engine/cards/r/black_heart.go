package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Black Heart (Dice Shard/Soul Heart Card)
//
//	Choose a player. Prevent the next 1 damage they would take this turn. Till end of turn, when that player dies, deal 2 damage to each player other than that player and you.
var blackHeart = engine.CardDef{
	Ref:    "black_heart",
	Kind:   engine.LootCard,
	Copies: 2,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Choose a player. Prevent the next 1 damage they would take this turn. Till end of turn, when that player dies, deal 2 damage to each player other than that player and you.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
		Effects: []engine.Effect{engine.PreventDamage(1, 0), engine.EffectFunc(func(c *engine.Ctx) {
			// The lingering trigger finds the chosen player by this mark.
			cardkit.MarkUsedOn(c.G, c.Source, int(c.Targets[0].Player)+1)
		}), engine.Linger(0)},
	}},
	Lingering: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When that player dies, deal 2 damage to each player other than that player and you.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Player != engine.NoPlayer && cardkit.UsedOn(g, self, int(e.Player)+1)
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { eachOtherThan(c, c.EventPlayer, 2) })},
	}},
}
