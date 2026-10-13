package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Distant Admiration (Active Treasure Card)
//
//	{Tap Effect}Choose another player. They give you a loot card. Reveal it, then you must play that loot card if able. This doesn't use a loot play.
var distantAdmiration = engine.CardDef{
	Ref:    "distant_admiration",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Choose another player. They give you a loot card. Reveal it, then you must play that loot card if able. This doesn't use a loot play.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
		Effects: []engine.Effect{giveLootThenPlay},
	}},
}
