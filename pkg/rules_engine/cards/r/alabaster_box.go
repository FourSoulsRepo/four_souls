package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Alabaster Box (Passive Treasure Card)
//
//	At the start of your turn, put a counter on this.
//	When the 2nd counter is put on this, gain 10¢.
//	When the 4th counter is put on this, loot 4.
//	When the 6th counter is put on this, gain +2 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var alabasterBox = engine.CardDef{
	Ref:    "alabaster_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
