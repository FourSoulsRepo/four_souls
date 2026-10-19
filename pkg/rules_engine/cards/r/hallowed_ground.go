package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hallowed Ground (Paid Treasure Card)
//
//	Each time a player rolls a ❷, put a counter on this.
//	{Paid Effect}Remove 3 counters from this:
//	Put the top card of the loot discard into your hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hallowedGround = engine.CardDef{
	Ref:    "hallowed_ground",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
