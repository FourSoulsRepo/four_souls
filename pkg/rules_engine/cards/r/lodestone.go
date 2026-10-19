package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lodestone (Passive Treasure Card)
//
//	Each time a player rolls a ❸, choose a monster. It gains -1{DC}, till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lodestone = engine.CardDef{
	Ref:    "lodestone",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
