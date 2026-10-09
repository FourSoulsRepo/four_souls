package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chaos Card (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose one- Kill a player or monster. Destroy an item or soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chaosCard = engine.CardDef{
	Ref:    "chaos_card",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
