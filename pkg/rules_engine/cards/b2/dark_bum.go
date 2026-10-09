package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark Bum (Passive Treasure Card)
//
//	At the start of your turn, roll-
//	1-2: Gain 3¢.
//	3-4: Loot 1.
//	5-6: Take 1 damage.
var darkBum = engine.CardDef{
	Ref:    "dark_bum",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, roll- 1-2: Gain 3¢. 3-4: Loot 1. 5-6: Take 1 damage.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainCents(3)).
				Results(3, 4, engine.Loot(1)).
				Results(5, 6, engine.DealDamage(1, engine.You)))},
		},
	},
}
