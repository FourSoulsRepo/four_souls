package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Reroll an item you control.
//	3-4: Reroll an item.
//	5-6: Reroll each item you control.
var pills2Fsp2 = engine.CardDef{
	Ref:    "pills_2_fsp2",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Roll- 1-2: Reroll an item you control. 3-4: Reroll an item. 5-6: Reroll each item you control.",
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, rerollAnItem(true)).
			Results(3, 4, rerollAnItem(false)).
			Results(5, 6, rerollYourItems))},
	}},
}
