package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Lust (Bonus Soul Card)
//
//	Each time a player kills a monster, put a counter on this.
//	When the 6th counter is put on this, the active player gains this soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulOfLust = engine.CardDef{
	Ref:    "soul_of_lust",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
}
