package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Head Of The Keeper (Passive Treasure Card)
//
//	Each time you deal damage, gain 1¢.
//	Each time you kill a monster or player, gain 2¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var headOfTheKeeper = engine.CardDef{
	Ref:    "head_of_the_keeper",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
