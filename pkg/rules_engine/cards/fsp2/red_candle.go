package fsp2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Red Candle (Active Treasure Card)
//
//	{Tap Effect}Before a dice is rolled, choose a number. Till the end of turn, each time that number is rolled, deal 1 damage to a monster or player.
var redCandle = engine.CardDef{
	Ref:    "red_candle",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Before a dice is rolled, choose a number. Till the end of turn, each time that number is rolled, deal 1 damage to a monster or player.",
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if a[0] < 0 {
				return
			}
			cardkit.MarkUsedOn(c.G, c.Source, a[0]+1) // the chosen number
			c.Do(engine.Linger(0))
		}, engine.Question{Text: "Choose a number", Options: func(*engine.Ctx, []int) []string { return []string{"1", "2", "3", "4", "5", "6"} }})},
	}},
	Lingering: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time the chosen number is rolled, deal 1 damage to a monster or player.",
		Trigger: engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return cardkit.UsedOn(g, self, e.Amount)
		}},
		Effects: []engine.Effect{damageAnything(1)},
	}},
}
