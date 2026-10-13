package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Rainbow Baby (Active Treasure Card)
//
//	{Tap Effect}Choose one- Each player loots 1. Each player gains 2¢. Each player takes 1 damage.
var rainbowBaby = engine.CardDef{
	Ref:    "rainbow_baby",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Choose one- Each player loots 1. Each player gains 2¢. Each player takes 1 damage.",
		Costs: []engine.Cost{engine.Tap()},
		Modes: []engine.Mode{
			{Text: "Each player loots 1.", Effects: []engine.Effect{engine.EachPlayer(engine.Loot(1))}},
			{Text: "Each player gains 2¢.", Effects: []engine.Effect{engine.EachPlayer(engine.GainCents(2))}},
			{Text: "Each player takes 1 damage.", Effects: []engine.Effect{engine.EachPlayer(engine.DealDamage(1, engine.You))}},
		},
	}},
}
