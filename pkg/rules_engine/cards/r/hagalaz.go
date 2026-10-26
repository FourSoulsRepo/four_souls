package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	hagalazPutThemBack = engine.NewKey("card.hagalaz.put_them_back_in_any")
)

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
			}, engine.Question{Text: cardkit.AskLookAtWhichDeck, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				return deckLabels(monsterOrRoomDecks(c.G))
			}}, engine.Question{
				Text: hagalazPutThemBack,
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
