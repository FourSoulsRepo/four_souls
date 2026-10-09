package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Empty Vessel (Passive Treasure Card)
//
//	When you have 0 loot cards in your hand, you have +1{ATK}.
//	While you have 0¢, you have +1 to your attack rolls.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var emptyVessel = engine.CardDef{
	Ref:    "empty_vessel",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
