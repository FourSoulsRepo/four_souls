package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed’s Gullet (Passive Treasure Card)
//
//	Each time you die, before paying penalties, gain 8¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var greedsGullet = engine.CardDef{
	Ref:    "greeds_gullet",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
