package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Fire Mind (Passive Treasure Card)
//
//	Combat damage you deal is doubled on attack rolls of 6.
//	Combat damage you take is doubled on attack rolls of 1.
var fireMind = engine.CardDef{
	Ref:    "fire_mind",
	Kind:   engine.TreasureCard,
	Copies: 1,
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if (hit && roll == 6) || (!hit && roll == 1) {
			return 2 * n
		}
		return n
	},
}
