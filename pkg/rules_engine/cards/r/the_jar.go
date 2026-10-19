package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Jar (Paid Treasure Card)
//
//	{Paid Effect}Pay 4¢:
//	Choose a player or monster. Prevent the next 1 damage they would take this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theJar = engine.CardDef{
	Ref:    "the_jar",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
