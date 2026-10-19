package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Tnt Barrel (Good Event Card)
//
//	Choose a player or monster, then roll-
//	1-4: You and that player or monster each takes 2 damage.
//	5-6: That player or monster take 4 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var tntBarrel = engine.CardDef{
	Ref:    "tnt_barrel",
	Kind:   engine.EventCard,
	Copies: 1,
}
