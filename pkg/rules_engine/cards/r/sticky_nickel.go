package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sticky Nickel (Nickel Card)
//
//	Roll-
//	1-3: Give this to the player to your left.
//	4-6: Gain 5¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var stickyNickel = engine.CardDef{
	Ref:    "sticky_nickel",
	Kind:   engine.LootCard,
	Copies: 1,
}
