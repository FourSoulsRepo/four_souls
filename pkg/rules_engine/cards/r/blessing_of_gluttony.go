package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of Gluttony (Room Card)
//
//	The active player loots +1 during their loot step.
var blessingOfGluttony = engine.CardDef{
	Ref:     "blessing_of_gluttony",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{activePlayerHas(engine.StatLootStep, 1)},
}
