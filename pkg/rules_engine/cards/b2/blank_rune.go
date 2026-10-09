package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blank Rune (Pill/Rune Card)
//
//	Roll-
//	1: Each player gains 1¢.
//	2: Each player loots 2.
//	3: Each player takes 3 damage.
//	4: Each player gains 4¢.
//	5: Each player loots 5.
//	6: Each player gains 6¢.
var blankRune = engine.CardDef{
	Ref:    "blank_rune",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Roll- 1: Each player gains 1¢. 2: Each player loots 2. 3: Each player takes 3 damage. 4: Each player gains 4¢. 5: Each player loots 5. 6: Each player gains 6¢.",
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 1, engine.EachPlayer(engine.GainCents(1))).
				Results(2, 2, engine.EachPlayer(engine.Loot(2))).
				Results(3, 3, engine.EachPlayer(engine.DealDamage(3, engine.You))).
				Results(4, 4, engine.EachPlayer(engine.GainCents(4))).
				Results(5, 5, engine.EachPlayer(engine.Loot(5))).
				Results(6, 6, engine.EachPlayer(engine.GainCents(6))))},
		},
	},
}
