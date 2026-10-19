package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Heavy Is The Head (Room Card)
//
//	Players who control the most items or tied for the most may only recharge one item during their recharge step.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var heavyIsTheHead = engine.CardDef{
	Ref:    "heavy_is_the_head",
	Kind:   engine.RoomCard,
	Copies: 1,
}
