package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Shovel (Active Treasure Card)
//
//	{Tap Effect}Put a non-event monster card in discard on top of the monster deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theShovel = engine.CardDef{
	Ref:    "the_shovel",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
