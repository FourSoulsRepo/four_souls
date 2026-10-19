package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Handicapped Placard (Passive Treasure Card)
//
//	Monster have -x{DC} on your turn, where x is the number of souls the player with the most souls controls minus the number of souls you control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var handicappedPlacard = engine.CardDef{
	Ref:    "handicapped_placard",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
