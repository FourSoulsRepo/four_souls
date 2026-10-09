package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lucky Foot (Active Treasure Card)
//
//	{Tap Effect}Add up to 2 to a non-attack roll.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var luckyFoot = engine.CardDef{
	Ref:    "lucky_foot",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
