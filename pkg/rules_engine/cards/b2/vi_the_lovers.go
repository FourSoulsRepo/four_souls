package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// VI. The Lovers (Wildcard Card)
//
//	Choose a player.
//	They gain +2{HP} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var viTheLovers = engine.CardDef{
	Ref:    "vi_the_lovers",
	Kind:   engine.LootCard,
	Copies: 1,
}
