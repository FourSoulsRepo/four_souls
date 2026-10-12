package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Bible (Active Treasure Card)
//
//	{Tap Effect}Choose one-
//	Destroy a curse.
//	Cancel your attack.
//	Kill a monster named Mom!, Mom's Heart!, or It Lives!
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theBible = engine.CardDef{
	Ref:    "the_bible",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
