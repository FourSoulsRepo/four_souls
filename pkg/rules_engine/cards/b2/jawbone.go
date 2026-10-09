package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Jawbone (Active Treasure Card)
//
//	{Tap Effect}Steal 3¢ from a player.
var jawbone = engine.CardDef{
	Ref:    "jawbone",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Steal 3¢ from a player.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
			Effects: []engine.Effect{engine.StealCents(3, 0)},
		},
	},
}
