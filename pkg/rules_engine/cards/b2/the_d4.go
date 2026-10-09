package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D4 (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose a player. They reroll each item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theD4 = engine.CardDef{
	Ref:    "the_d4",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
