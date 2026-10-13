package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Divorce Papers (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose another player. They give you half of their ¢ and loot cards rounded down, then gives you an item.
var divorcePapers = engine.CardDef{
	Ref:    "divorce_papers",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Destroy this. If you do, choose another player. They give you half of their ¢ and loot cards rounded down, then give you an item.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
		Effects: []engine.Effect{engine.DestroyThis(divorce)},
	}},
}
