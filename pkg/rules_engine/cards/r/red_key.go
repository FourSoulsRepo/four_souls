package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Key (Active Treasure Card)
//
//	{Tap Effect}Look at the top 3 cards of the room or monster deck. You may put one of those in a slot and the rest back. This can't be activated during an attack.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var redKey = engine.CardDef{
	Ref:    "red_key",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
