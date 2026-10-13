package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Polyphemus (Passive Treasure Card)
//
//	+2{ATK}
//	Monsters have +1{DC} on your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var polyphemus = engine.CardDef{
	Ref:    "polyphemus",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
