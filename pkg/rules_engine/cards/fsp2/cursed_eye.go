package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Eye (Passive Treasure Card)
//
//	Combat damage you deal on attack rolls of 6 is increased by 3.
//	When you roll an attack roll of 1, end your turn. Cancel everything that hasn't resolved.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedEye = engine.CardDef{
	Ref:    "cursed_eye",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
