package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Bomb! (Bomb Card)
//
//	Deal 2 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var megaBomb = engine.CardDef{
	Ref:    "mega_bomb",
	Kind:   engine.LootCard,
	Copies: 1,
}
