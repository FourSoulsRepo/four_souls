package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Consolation Prize (Passive Treasure Card)
//
//	At the start of your turn, loot x, where x is the number of souls the player with the most souls controls minus the number of souls you control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var consolationPrize = engine.CardDef{
	Ref:    "consolation_prize",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
