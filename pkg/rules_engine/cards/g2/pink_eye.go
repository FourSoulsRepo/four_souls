package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pink Eye (Trinket Card)
//
//	Each time you take combat damage from a monster or player on an attack roll of 1, deal 1 damage to them.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pinkEye = engine.CardDef{
	Ref:    "pink_eye",
	Kind:   engine.LootCard,
	Copies: 1,
}
