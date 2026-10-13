package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Divorce Papers (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose another player. They give you half of their ¢ and loot cards rounded down, then gives you an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var divorcePapers = engine.CardDef{
	Ref:    "divorce_papers",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
