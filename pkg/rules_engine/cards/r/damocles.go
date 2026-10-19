package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Damocles (Passive Treasure Card)
//
//	If you would gain any number of treasures, instead gain that many +1.
//	When you die, before paying penalties, destroy this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var damocles = engine.CardDef{
	Ref:    "damocles",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
