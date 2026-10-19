package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Emergency Contact (Wildcard Card)
//
//	Cancel an attack on a monster and put that monster card on the bottom of the monster deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var emergencyContact = engine.CardDef{
	Ref:    "emergency_contact",
	Kind:   engine.LootCard,
	Copies: 1,
}
