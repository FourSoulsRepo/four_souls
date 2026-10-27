package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	goldenTrinketPutAGold = engine.NewKey("card.golden_trinket.put_a_gold_counter_on")
)

// Golden Trinket (Trinket Card)
//
//	When this enters play put a gold counter on another non-eternal item you control.
//	This has the abilities of other items with gold counters on them.
//	-Trinket- This loot becomes an item under your control when it resolves.
var goldenTrinket = engine.CardDef{
	Ref:     "golden_trinket",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	GainsAbilities: func(g *engine.Game, self engine.ObjectID) []engine.CardRef {
		var out []engine.CardRef
		for _, id := range itemsOf(g, g.Object(self).Controller) {
			if id != self && g.Object(id).CountersOf(gold) > 0 {
				out = append(out, g.CardOf(id))
			}
		}
		return out
	},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this enters play put a gold counter on another non-eternal item you control.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ids := otherNonEternal(c); a[0] >= 0 && a[0] < len(ids) {
				c.G.AddNamedCounters(ids[a[0]], gold, 1)
			}
		}, engine.Question{Text: goldenTrinketPutAGold, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			return labelsOf(c.G, otherNonEternal(c))
		}})},
	}},
}

// gold names gold counters.
const gold = "gold"

// otherNonEternal lists the controller's non-eternal items but this one.
func otherNonEternal(c *engine.Ctx) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range nonEternalItems(c.G, c.Controller) {
		if id != c.Source {
			out = append(out, id)
		}
	}
	return out
}
