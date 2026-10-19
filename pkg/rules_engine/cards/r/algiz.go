package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Algiz (Pill/Rune Card)
//
//	Choose a player. Each item they control gains eternal till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var algiz = engine.CardDef{
	Ref:    "algiz",
	Kind:   engine.LootCard,
	Copies: 1,
}
