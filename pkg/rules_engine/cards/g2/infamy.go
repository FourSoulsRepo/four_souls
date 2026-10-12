package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Infamy (Passive Treasure Card)
//
//	You take no combat damage on attack rolls of 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var infamy = engine.CardDef{
	Ref:    "infamy",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
