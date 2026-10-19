package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lemegeton (Eternal Treasure Card)
//
//	{Tap Effect}Gain +1 treasure. The next time your turn ends, destroy a non-eternal item you control.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lemegeton = engine.CardDef{
	Ref:     "lemegeton",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
