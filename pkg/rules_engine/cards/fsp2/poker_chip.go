package fsp2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Poker Chip (Trinket Card)
//
//	The first time you would gain ¢ on each of your turns, roll-
//	1-3: Gain 1 ¢ instead.
//	4-6: Gain double the number of ¢ you would've gained.
//	-Trinket- This loot becomes an item under your control when it resolves.
var pokerChip = engine.CardDef{
	Ref:     "poker_chip",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Replacements: []engine.Replacement{{
		Text: "The first time you would gain ¢ on each of your turns, roll- 1-3: Gain 1¢ instead. 4-6: Gain double the number of ¢ you would've gained.",
		When: func(g *engine.Game, self engine.ObjectID, a engine.Action) bool {
			p := g.Object(self).Controller
			return a.Kind == engine.ActGainCents && a.Amount > 0 && a.Player == p && g.Turn.Active == p && !cardkit.Used(g, self)
		},
		Do: func(g *engine.Game, self engine.ObjectID, a engine.Action) []engine.Action {
			cardkit.MarkUsed(g, self)
			return []engine.Action{engine.PushLingering(self, 0, a.Player, a.Amount)}
		},
	}},
	Lingering: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Roll- 1-3: Gain 1¢ instead. 4-6: Gain double the number of ¢ you would've gained.",
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 3, engine.GainCents(1)).
			Results(4, 6, engine.EffectFunc(func(c *engine.Ctx) { c.Do(engine.GainCents(2 * c.EventAmount)) })))},
	}},
}
