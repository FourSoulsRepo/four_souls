package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper’s Box (Passive Treasure Card)
//
//	When this enters play, expand shop slots by 2.
//	Each time another player purchases a shop item, gain 2¢ and loot 1.
var keepersBox = engine.CardDef{
	Ref:    "keepers_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this enters play, expand shop slots by 2.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandShop(2) })},
		},
		{
			Kind: engine.Triggered,
			Text: "Each time another player purchases a shop item, gain 2¢ and loot 1.",
			Trigger: engine.Trigger{On: engine.EvPurchased, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Text == "" && e.Player != g.Object(self).Controller
			}},
			Effects: []engine.Effect{engine.GainCents(2), engine.Loot(1)},
		},
	},
}
