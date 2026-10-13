package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mutant Spider (Active Treasure Card)
//
//	{Tap Effect}The next time a player would roll a dice, they instead roll 4 dice. You choose one of the rolls as the result.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mutantSpider = engine.CardDef{
	Ref:    "mutant_spider",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
