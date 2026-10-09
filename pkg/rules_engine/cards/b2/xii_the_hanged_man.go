package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XII. The Hanged Man (Wildcard Card)
//
//	Look at the top card of each deck. You may put any of those cards on the bottom of their deck, then loot 2.
var xiiTheHangedMan = engine.CardDef{
	Ref:    "xii_the_hanged_man",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Look at the top card of each deck. You may put any of those cards on the bottom of their deck, then loot 2.",
			Effects: []engine.Effect{engine.Ask(hangedMan, hangedManQuestions()...), engine.Loot(2)},
		},
	},
}

// hangedManQuestions asks, deck by deck, whether to keep the top card.
func hangedManQuestions() []engine.Question {
	var qs []engine.Question
	for k := range 4 { // at most 4 decks
		qs = append(qs, engine.Question{
			Text: "Put this card on the bottom of its deck?",
			Options: func(c *engine.Ctx, _ []int) []string {
				decks := c.G.DecksInPlay()
				if k >= len(decks) || len(c.G.Decks[decks[k]]) == 0 {
					return nil
				}
				top := c.G.Object(c.G.DeckTop(decks[k], 1)[0]).Card
				return []string{"keep " + string(top) + " on top", "put " + string(top) + " on the bottom"}
			},
		})
	}
	return qs
}

func hangedMan(c *engine.Ctx, a []int) {
	for k, d := range c.G.DecksInPlay() {
		top := c.G.DeckTop(d, 1)
		if len(top) == 0 {
			continue
		}
		c.G.LookAt(c.Controller, top...)
		if a[k] == 1 {
			c.G.DeckToBottom(d, top[0])
		}
	}
}
