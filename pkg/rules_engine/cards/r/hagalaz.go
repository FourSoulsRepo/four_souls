package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hagalaz (Pill/Rune Card)
//
//	Look at the top 5 cards of the monster or room deck and put them back in any order, then put a room or monster not being attacked into discard.
var hagalaz = engine.CardDef{
	Ref:    "hagalaz",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Look at the top 5 cards of the monster or room deck and put them back in any order, then put a room or monster not being attacked into discard.",
		Effects: []engine.Effect{
			engine.Ask(func(c *engine.Ctx, a []int) {
				if ds := monsterOrRoomDecks(c.G); a[0] >= 0 && a[0] < len(ds) {
					c.G.SetDeckTop(ds[a[0]], c.Order)
				}
			}, engine.Question{Text: "Look at which deck?", Options: func(c *engine.Ctx, _ []int) []string {
				return deckLabels(monsterOrRoomDecks(c.G))
			}}, engine.Question{
				Text: "Put them back in any order, top first.",
				Order: func(c *engine.Ctx, a []int) []engine.ObjectID {
					if ds := monsterOrRoomDecks(c.G); a[0] >= 0 && a[0] < len(ds) {
						return c.G.DeckTop(ds[a[0]], 5)
					}
					return nil
				},
			}),
			roomOrMonsterToDiscard(false),
		},
	}},
}
