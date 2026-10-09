package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// There’s Options (Passive Treasure Card)
//
//	You may look at the top card of the treasure deck at any time on your turn.
//	You may purchase an additional time on your turn.
var theresOptions = engine.CardDef{
	Ref:           "theres_options",
	Kind:          engine.TreasureCard,
	Copies:        1,
	PeeksTreasure: true,
	Statics:       []engine.Static{engine.YouHave(engine.StatPurchases, 1)},
}
