package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Magnet (Passive Treasure Card)
//
//	Each time another player gains ¢, they must give you 1¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var magnet = engine.CardDef{
	Ref:    "magnet",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
