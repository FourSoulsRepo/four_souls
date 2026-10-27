package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	lemegetonDestroyWhichNon = engine.NewKey("card.lemegeton.destroy_which_non_eternal_item")
)

// Lemegeton (Eternal Treasure Card)
//
//	{Tap Effect}Gain +1 treasure. The next time your turn ends, destroy a non-eternal item you control.
//	-Eternal- This can't be destroyed or put into discard.
var lemegeton = engine.CardDef{
	Ref:     "lemegeton",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Gain +1 treasure. The next time your turn ends, destroy a non-eternal item you control.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.GainTreasure(1), engine.EffectFunc(func(c *engine.Ctx) {
				c.G.AddNamedCounters(c.Source, lemegetonDue, 1)
			})},
		},
		{
			Kind: engine.Triggered,
			Text: "The next time your turn ends, destroy a non-eternal item you control.",
			Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				o := g.Object(self)
				return e.Player == o.Controller && o.CountersOf(lemegetonDue) > 0
			}},
			Effects: []engine.Effect{lemegetonDestroys()},
		},
	},
}

// lemegetonDue counts the destructions owed at the end of the turn.
const lemegetonDue = "due"

// lemegetonDestroys: one item per use; each is chosen alone, so one item
// may be chosen each time and the rest do nothing (Q-15, rules as
// written).
func lemegetonDestroys() engine.Effect {
	qs := make([]engine.Question, 6)
	for i := range qs {
		qs[i] = engine.Question{Text: lemegetonDestroyWhichNon, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			if i >= c.G.Object(c.Source).CountersOf(lemegetonDue) {
				return nil
			}
			return labelsOf(c.G, nonEternalItems(c.G, c.Controller))
		}}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		items := nonEternalItems(c.G, c.Controller)
		c.G.AddNamedCounters(c.Source, lemegetonDue, -c.G.Object(c.Source).CountersOf(lemegetonDue))
		for _, x := range a {
			if x >= 0 && x < len(items) {
				c.G.DestroyObject(c.Controller, items[x])
			}
		}
	}, qs...)
}
