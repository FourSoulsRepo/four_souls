package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Betrayal! (Good Event Card)
//
//	Choose another player. Roll-
//	1-2: Steal 5¢ from them.
//	3-4: Steal a loot card from them at random.
//	5: Steal a non-eternal item they control.
//	6: Steal a soul they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var betrayal = engine.CardDef{
	Ref:    "betrayal",
	Kind:   engine.EventCard,
	Copies: 1,
	Soul:   1,
}
