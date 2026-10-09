package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Monstro’s Tooth (Passive Treasure Card)
//
//	At the start of your turn, choose a player at random. That player destroys an item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var monstrosTooth = engine.CardDef{
	Ref:    "monstros_tooth",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
