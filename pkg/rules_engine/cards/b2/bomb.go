package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bomb! (Bomb Card)
//
//	Deal 1 damage to a Monster or Player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bomb = engine.CardDef{
	Ref:    "bomb",
	Kind:   engine.LootCard,
	Copies: 4,
}
