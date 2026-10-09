package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// III. The Empress (Wildcard Card)
//
//	Choose a player.
//	They gain +1{ATK} and +1 to dice rolls till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var iiiTheEmpress = engine.CardDef{
	Ref:    "iii_the_empress",
	Kind:   engine.LootCard,
	Copies: 1,
}
