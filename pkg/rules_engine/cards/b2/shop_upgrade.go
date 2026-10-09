package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Shop Upgrade! (Good Event Card)
//
//	Expand shop slots by 2.
//	The active player may attack an additional time this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var shopUpgrade = engine.CardDef{
	Ref:    "shop_upgrade",
	Kind:   engine.EventCard,
	Copies: 1,
}
