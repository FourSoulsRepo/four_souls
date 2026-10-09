package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Bra (Active Treasure Card)
//
//	{Tap Effect}Choose a monster or player. The next instance of damage they take this turn is reduced to 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsBra = engine.CardDef{
	Ref:    "moms_bra",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
