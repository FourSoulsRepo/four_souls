package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Suicide King (Passive Treasure Card)
//
//	Each time you die, before paying penalties, loot 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var suicideKing = engine.CardDef{
	Ref:    "suicide_king",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
