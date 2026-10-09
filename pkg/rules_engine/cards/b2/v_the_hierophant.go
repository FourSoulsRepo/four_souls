package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// V. The Hierophant (Wildcard Card)
//
//	Choose a player or monster. Prevent the next instance of up to 2 damage they would take this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var vTheHierophant = engine.CardDef{
	Ref:    "v_the_hierophant",
	Kind:   engine.LootCard,
	Copies: 1,
}
