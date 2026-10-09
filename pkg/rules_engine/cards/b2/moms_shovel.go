package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Shovel (One-Use Treasure Card)
//
//	This enters play deactivated.
//	{Tap Effect}Destroy this. If you do, steal a soul from another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsShovel = engine.CardDef{
	Ref:    "moms_shovel",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
	Tap:    true,
}
