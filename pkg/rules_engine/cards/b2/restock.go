package b2

import engine "github.com/FourSoulsRepo/rules_engine"

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
		qs = append(qs, engine.Question{Text: "Discard this shop item?", Options: func(c *engine.Ctx, _ []int) []string {
			items := c.G.ShopItems()
			if k >= len(items) {
				return nil
			}
			card := string(c.G.Object(items[k]).Card)
			return []string{"keep " + card, "discard " + card}
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
