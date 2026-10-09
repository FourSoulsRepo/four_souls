package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mr. Boom (Active Treasure Card)
//
//	{Tap Effect}Deal 1 damage to a monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mrBoom = engine.CardDef{
	Ref:    "mr_boom",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
