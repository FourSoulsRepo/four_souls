package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pay To Play (Paid Treasure Card)
//
//	{Paid Effect}Pay 10¢:
//	Steal a non-eternal item a player controls.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var payToPlay = engine.CardDef{
	Ref:    "pay_to_play",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
