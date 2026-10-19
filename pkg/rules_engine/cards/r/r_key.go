package r

import engine "github.com/FourSoulsRepo/rules_engine"

// R Key (One-Use Treasure Card)
//
//	This enters play deactivated.
//	{Tap Effect}Destroy this. If you do, destroy all souls. Each player discards their hand and loots 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var rKey = engine.CardDef{
	Ref:    "r_key",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
