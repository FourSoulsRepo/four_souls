package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Loot 1.
//	3-4: Loot 3.
//	5-6: Discard 1 loot card.
var pills = engine.CardDef{
	Ref:    "pills",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Roll- 1-2: Loot 1. 3-4: Loot 3. 5-6: Discard 1 loot card.",
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.Loot(1)).
				Results(3, 4, engine.Loot(3)).
				Results(5, 6, discardOne))},
		},
	},
}
