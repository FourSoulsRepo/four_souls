package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Wiz (Passive Treasure Card)
//
//	Monsters have +1{DC} on your turn.
//	Each time you deal combat damage, deal 1 damage to another monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theWiz = engine.CardDef{
	Ref:    "the_wiz",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
