package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dry Baby (Passive Treasure Card)
//
//	Damage you would take is reduced to 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dryBaby = engine.CardDef{
	Ref:    "dry_baby",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
