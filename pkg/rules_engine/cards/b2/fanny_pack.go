package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fanny Pack (Passive Treasure Card)
//
//	Each time you take damage, loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fannyPack = engine.CardDef{
	Ref:    "fanny_pack",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
