package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	sacredHeartChangeThe1 = engine.NewKey("card.sacred_heart.change_the_1_to_a")
)

// Sacred Heart (Passive Treasure Card)
//
//	When you would roll a 1, you may change the result to a 6.
var sacredHeart = engine.CardDef{
	Ref:    "sacred_heart",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "When you would roll a 1, you may change the result to a 6.",
			Trigger: engine.Trigger{On: engine.EvRollWouldResolve, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Amount == 1 && e.Player == g.Object(self).Controller
			}},
			Effects: []engine.Effect{may(sacredHeartChangeThe1, engine.EffectFunc(func(c *engine.Ctx) { c.G.SetRoll(c.EventStack, 6, c.Source) }))},
		},
	},
}
