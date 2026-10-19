package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pity For The Poor (Room Card)
//
//	Players who control the fewest souls or tied for fewest may purchase a shop item for 0¢ on their turn.
var pityForThePoor = engine.CardDef{
	Ref:    "pity_for_the_poor",
	Kind:   engine.RoomCard,
	Copies: 1,
	// Every purchase is free while they have the fewest souls (Q-21).
	Statics: []engine.Static{playersWith(engine.StatShopPrice, -1000, fewestOf(soulCount))},
}
