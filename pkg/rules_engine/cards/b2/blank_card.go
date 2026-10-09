package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blank Card (Active Treasure Card)
//
//	{Tap Effect}The next time you play a non-trinket, non-ambush loot card this turn, copy it.
var blankCard = engine.CardDef{
	Ref:    "blank_card",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: The next time you play a non-trinket, non-ambush loot card this turn, copy it.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.Players[c.Controller].CopyNextLoot = true })},
		},
	},
}
