package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Shard (Dice Shard/Soul Heart Card)
//
//	Choose a dice roll. Its controller rerolls it.
var diceShard = engine.CardDef{
	Ref:    "dice_shard",
	Kind:   engine.LootCard,
	Copies: 3,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a dice roll. Its controller rerolls it.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.RerollRoll(0)},
		},
	},
}
