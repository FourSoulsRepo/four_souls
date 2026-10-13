package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Candle (Active Treasure Card)
//
//	{Tap Effect}Before a dice is rolled, choose a number. Till the end of turn, each time that number is rolled, deal 1 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var redCandle = engine.CardDef{
	Ref:    "red_candle",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
