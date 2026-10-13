package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Euthanasia (Passive Treasure Card)
//
//	Each time you roll the same result twice in a row on an attack roll on the same turn, kill the monster you're attacking.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var euthanasia = engine.CardDef{
	Ref:    "euthanasia",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
