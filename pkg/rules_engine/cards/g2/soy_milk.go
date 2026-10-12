package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soy Milk (Passive Treasure Card)
//
//	Monsters have +1{DC} on your turn.
//	Each time you deal combat damage to a monster, deal 1 damage to each other monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soyMilk = engine.CardDef{
	Ref:    "soy_milk",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
