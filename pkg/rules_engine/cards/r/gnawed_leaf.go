package r

import (
	"slices"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Gnawed Leaf (Passive Treasure Card)
//
//	Prevent all damage you would take while it's not your turn.
//	The first time each turn another player plays a loot card that targets you or something you control, you may cancel it.
var gnawedLeaf = engine.CardDef{
	Ref:       "gnawed_leaf",
	Kind:      engine.TreasureCard,
	Copies:    1,
	ThreePlus: true,
	DamageMod: func(g *engine.Game, self engine.ObjectID, t engine.Target, n int, _ bool) int {
		if o := g.Object(self); t.IsPlayer && t.Player == o.Controller && g.Turn.Active != o.Controller && o.Role == engine.RoleItem {
			return 0
		}
		return n
	},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "The first time each turn another player plays a loot card that targets you or something you control, you may cancel it.",
		Trigger: engine.Trigger{On: engine.EvStackAdded, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			me := g.Object(self).Controller
			if cardkit.Used(g, self) {
				return false
			}
			for _, it := range g.Stack {
				if it.Kind == engine.StackLoot && it.Source == e.Object && it.Controller != me {
					return slices.ContainsFunc(it.Targets, func(t engine.Chosen) bool {
						return (t.Kind == engine.TargetPlayer && t.Player == me) || (t.Object != 0 && g.Object(t.Object).Controller == me)
					})
				}
			}
			return false
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			cardkit.MarkUsed(c.G, c.Source)
		}), cardkit.May("Cancel that loot card?", engine.EffectFunc(func(c *engine.Ctx) {
			for _, it := range c.G.Stack {
				if it.Kind == engine.StackLoot && it.Source == c.EventObject {
					c.G.CancelStackItem(it.ID)
					return
				}
			}
		}))},
	}},
}
