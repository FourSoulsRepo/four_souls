package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// 9 Volt (Passive Treasure Card)
//
//	At the end of your turn, you may recharge up to 2 items you control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var card9Volt = engine.CardDef{
	Ref:    "9_volt",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
