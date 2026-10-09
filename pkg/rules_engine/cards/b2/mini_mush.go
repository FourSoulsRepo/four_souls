package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mini Mush (Active Treasure Card)
//
//	{Tap Effect}Subtract up to 2 from a roll.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var miniMush = engine.CardDef{
	Ref:    "mini_mush",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
