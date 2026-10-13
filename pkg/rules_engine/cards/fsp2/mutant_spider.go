package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mutant Spider (Active Treasure Card)
//
//	{Tap Effect}The next time a player would roll a dice, they instead roll 4 dice. You choose one of the rolls as the result.
var mutantSpider = engine.CardDef{
	Ref:    "mutant_spider",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: The next time a player would roll a dice, they instead roll 4 dice. You choose one of the rolls as the result.",
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.ManyDice = append(c.G.ManyDice, engine.ManyDice{Chooser: c.Controller, Dice: 4}) // R-DICE-11
		})},
	}},
}
