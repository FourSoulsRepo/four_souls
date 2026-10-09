package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curved Horn (Trinket Card)
//
//	Gain +1{ATK} for your first attack roll each turn.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curvedHorn = engine.CardDef{
	Ref:    "curved_horn",
	Kind:   engine.LootCard,
	Copies: 1,
}
