package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gold Bomb!! (Bomb Card)
//
//	Deal 3 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var goldBomb = engine.CardDef{
	Ref:    "gold_bomb",
	Kind:   engine.LootCard,
	Copies: 1,
}
