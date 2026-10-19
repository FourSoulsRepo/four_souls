package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Backstabber (Passive Treasure Card)
//
//	You may attack players who control more souls than you. They have +3{DC} for the attack.
//	Each time you kill a player, loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var backstabber = engine.CardDef{
	Ref:    "backstabber",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
