package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Berkano (Pill/Rune Card)
//
//	Destroy an item. If that item was controlled by a player, they steal an item from the shop.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var berkano = engine.CardDef{
	Ref:    "berkano",
	Kind:   engine.LootCard,
	Copies: 1,
}
