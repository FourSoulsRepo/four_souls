package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	shadowOfWarDestroyWhichItem = engine.NewKey("card.shadow_of_war.destroy_which_item")
)

// Shadow Of War (Room Card)
//
//	When a player dies, before paying penalties, they must destroy an item they control.
var shadowOfWar = engine.CardDef{
	Ref:    "shadow_of_war",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When a player dies, before paying penalties, they must destroy an item they control.",
		Takes:   engine.TakesEventPlayer,
		Trigger: engine.WhenAPlayerDies(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ids := nonEternalItems(c.G, c.EventPlayer); a[0] >= 0 && a[0] < len(ids) {
				c.G.DestroyObject(c.EventPlayer, ids[a[0]])
			}
		}, engine.Question{
			Text:    shadowOfWarDestroyWhichItem,
			Player:  func(c *engine.Ctx, _ []int) engine.PlayerID { return c.EventPlayer },
			Options: func(c *engine.Ctx, _ []int) []engine.Msg { return labelsOf(c.G, nonEternalItems(c.G, c.EventPlayer)) },
		})},
	}},
}
