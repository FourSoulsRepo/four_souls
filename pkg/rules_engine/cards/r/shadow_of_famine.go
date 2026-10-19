package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Shadow Of Famine (Room Card)
//
//	Players can't activate more than one ability each turn.
//	Players can't play more than one loot card each turn.
var shadowOfFamine = engine.CardDef{
	Ref:    "shadow_of_famine",
	Kind:   engine.RoomCard,
	Copies: 1,
	Statics: []engine.Static{
		playersHave(engine.StatOneAbility, 1),
		playersHave(engine.StatOneLoot, 1),
	},
}
