package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Halo Of Flies (Passive Treasure Card)
//
//	You take no combat damage on attack rolls of 1 or 6.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var haloOfFlies = engine.CardDef{
	Ref:    "halo_of_flies",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
