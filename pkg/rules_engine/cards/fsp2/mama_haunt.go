package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mama Haunt (Passive Treasure Card)
//
//	{Curse Effect}Your character doesn't recharge during your recharge step.
//	When you die, before paying penalties, give this to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mamaHaunt = engine.CardDef{
	Ref:    "mama_haunt",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
