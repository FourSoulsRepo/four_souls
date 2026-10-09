package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Shovel (Active Treasure Card)
//
//	{Tap Effect}Put a non-event monster card in discard on top of the monster deck.
var theShovel = engine.CardDef{
	Ref:    "the_shovel",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Put a non-event monster card in discard on top of the monster deck.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					c.G.DiscardToDeckTop(engine.MonsterDeck, nonEventDiscards(c.G)[a[0]])
				}
			}, engine.Question{Text: "Which monster card?", Options: func(c *engine.Ctx, _ []int) []string {
				var out []string
				for _, id := range nonEventDiscards(c.G) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			}})},
		},
	},
}

func nonEventDiscards(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Discards[engine.MonsterDeck] {
		if g.Kind(id) == engine.MonsterCard {
			out = append(out, id)
		}
	}
	return out
}
