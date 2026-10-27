package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	friendlySackPutTheRest   = engine.NewKey("card.friendly_sack.put_the_rest_in_which")
	friendlySackPutWhichCard = engine.NewKey("card.friendly_sack.put_which_card_in_your")
)

// Friendly Sack (Active Treasure Card)
//
//	{Tap Effect}Look at the top 3 cards of the loot deck. Put one in your hand and put the rest in another player's hand.
var friendlySack = engine.CardDef{
	Ref:    "friendly_sack",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Look at the top 3 cards of the loot deck. Put one in your hand and put the rest in another player's hand.",
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.LookAt(c.Controller, c.G.DeckTop(engine.LootDeck, 3)...)
		}), engine.Ask(func(c *engine.Ctx, a []int) {
			top := c.G.DeckTop(engine.LootDeck, 3)
			if a[0] < 0 || a[0] >= len(top) {
				return
			}
			c.G.DeckToHand(engine.LootDeck, top[a[0]], c.Controller)
			others := otherPlayers(c)
			if a[1] < 0 || a[1] >= len(others) {
				return
			}
			for i, id := range top {
				if i != a[0] {
					c.G.DeckToHand(engine.LootDeck, id, others[a[1]])
				}
			}
		}, engine.Question{Text: friendlySackPutWhichCard, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			return labelsOf(c.G, c.G.DeckTop(engine.LootDeck, 3))
		}}, engine.Question{Text: friendlySackPutTheRest, Options: func(c *engine.Ctx, a []int) []engine.Msg {
			if a[0] < 0 || len(c.G.DeckTop(engine.LootDeck, 3)) < 2 {
				return nil
			}
			return playerLabels(c.G, otherPlayers(c))
		}})},
	}},
}
