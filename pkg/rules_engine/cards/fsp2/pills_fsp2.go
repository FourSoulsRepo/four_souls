package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Cancel the effect of a loot being played.
//	3-4: Each other player discards a loot card.
var pillsFsp2 = engine.CardDef{
	Ref:    "pills_fsp2",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.LootAbility,
		Text:  "Roll- 1-2: Cancel the effect of a loot being played. 3-4: Each other player discards a loot card.",
		Takes: engine.TakesOthers,
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, cancelALoot).
			Results(3, 4, engine.EachOtherPlayer(discardOne)))},
	}},
}
