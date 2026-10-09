package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bum-Bo! (Passive Treasure Card)
//
//	If you would gain any amount of ¢, this levels up by that much instead.
//	{LV1 Effect}You have +2 to your first attack roll each turn.
//	{LV10 Effect}You have +1{ATK}.
//	{LV25 Effect}You may attack any number of times on your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bumBo = engine.CardDef{
	Ref:    "bum_bo",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
