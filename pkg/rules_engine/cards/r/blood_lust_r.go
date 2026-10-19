package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Blood Lust (Room Card)
//
//	The active player must attack each turn if able.
//	The first time the active player declares an attack each turn, they must attack an additional time this turn.
var bloodLustR = engine.CardDef{
	Ref:    "blood_lust_r",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "The active player must attack each turn if able.",
			Trigger: atStartOfTheTurn(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.Turn.MustAttacks++ })},
		},
		{
			Kind: engine.Triggered,
			Text: "The first time the active player declares an attack each turn, they must attack an additional time this turn.",
			Trigger: engine.Trigger{On: engine.EvAttackDeclared, Match: func(g *engine.Game, self engine.ObjectID, _ engine.Event) bool {
				return !cardkit.Used(g, self)
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if !cardkit.Used(c.G, c.Source) {
					cardkit.MarkUsed(c.G, c.Source)
					c.G.ForceAttacks(1, false)
				}
			})},
		},
	},
}
