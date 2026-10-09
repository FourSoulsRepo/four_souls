package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Butter Bean! (Butter Bean Card)
//
//	Cancel the ↷ or $ ability of an item or a loot being played.
var butterBean = engine.CardDef{
	Ref:    "butter_bean",
	Kind:   engine.LootCard,
	Copies: 3,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Cancel the ↷ or $ ability of an item or a loot being played.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetStackAbility)},
			Effects: []engine.Effect{engine.CancelTarget(0)},
		},
	},
}
