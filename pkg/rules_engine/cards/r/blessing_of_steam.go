package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of Steam (Room Card)
//
//	Shop items the active player purchases cost 5¢ less.
var blessingOfSteam = engine.CardDef{
	Ref:     "blessing_of_steam",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{activePlayerHas(engine.StatShopPrice, -5)},
}
