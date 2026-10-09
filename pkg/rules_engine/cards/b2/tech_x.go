package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Tech X (Active Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove 3 counters from this:
//	Kill a player or monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var techX = engine.CardDef{
	Ref:    "tech_x",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
