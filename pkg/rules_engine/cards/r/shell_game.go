package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Shell Game (Room Card)
//
//	{Paid Effect}Pay 3¢: roll-
//	3: Gain 1¢.
//	4: Loot 1.
//	5: Gain 5¢.
//	6: Gain +1 treasure. Put this into discard.
var shellGame = engine.CardDef{
	Ref:    "shell_game",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "$ Pay 3¢: Roll- 3: Gain 1¢. 4: Loot 1. 5: Gain 5¢. 6: Gain +1 treasure. Put this into discard.",
		Costs: []engine.Cost{engine.PayCents(3)},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(3, 3, engine.GainCents(1)).
			Results(4, 4, engine.Loot(1)).
			Results(5, 5, engine.GainCents(5)).
			Results(6, 6, engine.GainTreasure(1), discardThisRoom))},
	}},
}
