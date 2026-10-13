package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Daddy Long Legs (Passive Treasure Card)
//
//	At the start of your turn, deal 1 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var daddyLongLegsFsp2 = engine.CardDef{
	Ref:    "daddy_long_legs_fsp2",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
