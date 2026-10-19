package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mama Mega (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, deal 3 damage to each monster and player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mamaMega = engine.CardDef{
	Ref:    "mama_mega",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
