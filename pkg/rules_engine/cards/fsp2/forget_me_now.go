package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Forget Me Now (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, each player destroys a soul they control.
var forgetMeNow = engine.CardDef{
	Ref:    "forget_me_now",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Destroy this. If you do, each player destroys a soul they control.",
		Takes:   engine.TakesOthers,
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.DestroyThis(eachPlayerDestroysASoul)},
	}},
}
