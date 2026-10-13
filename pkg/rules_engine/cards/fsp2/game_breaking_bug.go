package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Game Breaking Bug! (Passive Treasure Card)
//
//	Each time a player rolls a ❶, you may reroll an item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gameBreakingBug = engine.CardDef{
	Ref:    "game_breaking_bug",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
