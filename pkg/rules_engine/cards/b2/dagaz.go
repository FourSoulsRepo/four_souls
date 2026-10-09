package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dagaz (Pill/Rune Card)
//
//	Choose one- Destroy a curse. Choose a player. Prevent the next 1 damage they would take this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dagaz = engine.CardDef{
	Ref:    "dagaz",
	Kind:   engine.LootCard,
	Copies: 1,
}
