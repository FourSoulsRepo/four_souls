package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Steamy Sale! (Passive Treasure Card)
//
//	Shop items you purchase cost 5¢ less.
var steamySale = engine.CardDef{
	Ref:     "steamy_sale",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatShopPrice, -5)},
}
