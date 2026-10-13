package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Boss Rush! (Good Event Card)
//
//	Reveal cards from the top of the monster deck till you reveal 2 boss cards. Put them in one or more monster slots not being attacked and the rest into discard. The active player must make an additional attack on one of them this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bossRush = engine.CardDef{
	Ref:    "boss_rush",
	Kind:   engine.EventCard,
	Copies: 1,
}
