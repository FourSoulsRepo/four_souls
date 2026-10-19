package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Eden’s Blessing (Room Card)
//
//	At the start of the turn, the active player may reroll an item they control.
var edensBlessingR = engine.CardDef{
	Ref:    "edens_blessing_r",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of the turn, the active player may reroll an item they control.",
		Trigger: atStartOfTheTurn(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ids := itemsOf(c.G, c.Controller); a[0] >= 0 && a[0] < len(ids) {
				c.G.RerollItem(ids[a[0]])
			}
		}, engine.Question{Text: "Reroll which item?", Options: func(c *engine.Ctx, _ []int) []string {
			return append(labelsOf(c.G, itemsOf(c.G, c.Controller)), "none")
		}})},
	}},
}
