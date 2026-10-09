package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lazarus’ Rags (Eternal Treasure Card)
//
//	Each time you die, after paying penalties, gain +1 treasure.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lazarusRags = engine.CardDef{
	Ref:     "lazarus_rags",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
