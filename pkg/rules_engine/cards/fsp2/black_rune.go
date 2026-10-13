package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Black Rune (Pill/Rune Card)
//
//	Roll-
//	1-2: Deal 1 damage to each monster.
//	3-4: Reroll an item.
//	5-6: Discard your hand, then loot 3.
var blackRune = engine.CardDef{
	Ref:    "black_rune",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Roll- 1-2: Deal 1 damage to each monster. 3-4: Reroll an item. 5-6: Discard your hand, then loot 3.",
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, engine.EachMonsterTakesDamage(1)).
			Results(3, 4, rerollAnItem(false)).
			Results(5, 6, discardHandThenLoot(3)))},
	}},
}
