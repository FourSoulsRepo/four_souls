package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Black Candle (Active Treasure Card)
//
//	{Tap Effect}Reveal the top 6 cards of the monster deck. Give any curse cards revealed to the player or players of your choosing. Put the rest on the bottom of the deck in any order.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blackCandle = engine.CardDef{
	Ref:    "black_candle",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
