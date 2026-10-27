package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	restockDiscardItem     = engine.NewKey("card.restock.discard")
	restockKeepItem        = engine.NewKey("card.restock.keep")
	restockDiscardThisShop = engine.NewKey("card.restock.discard_this_shop_item")
)

// Restock (Passive Treasure Card)
//
//	At the start of your turn, you may put any number of shop items into discard.
var restock = engine.CardDef{
	Ref:    "restock",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, you may put any number of shop items into discard.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.Ask(restockDiscard, restockQuestions()...)},
		},
	},
}

// restockQuestions ask about each shop slot in turn.
func restockQuestions() []engine.Question {
	var qs []engine.Question
	for k := range 6 {
		qs = append(qs, engine.Question{Text: restockDiscardThisShop, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			items := c.G.ShopItems()
			if k >= len(items) {
				return nil
			}
			card := engine.CardArg("card", c.G.Object(items[k]).Card)
			return []engine.Msg{restockKeepItem.Msg(card), restockDiscardItem.Msg(card)}
		}})
	}
	return qs
}

func restockDiscard(c *engine.Ctx, a []int) {
	items := c.G.ShopItems()
	for k, id := range items {
		if k < len(a) && a[k] == 1 {
			c.G.DiscardShopItem(id)
		}
	}
}
