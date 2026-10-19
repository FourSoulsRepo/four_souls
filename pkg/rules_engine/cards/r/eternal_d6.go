package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Eternal D6 (Active Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove 2 counters from this:
//	Till end of turn, if a player would roll a 2 or 3, it becomes a 1 instead.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var eternalD6 = engine.CardDef{
	Ref:    "eternal_d6",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
