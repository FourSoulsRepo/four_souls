package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Perthro (Pill/Rune Card)
//
//	Reroll an item.
var perthro = engine.CardDef{
	Ref:    "perthro",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Reroll an item.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetItem)},
		Effects: []engine.Effect{engine.RerollTarget(0)},
	}},
}
