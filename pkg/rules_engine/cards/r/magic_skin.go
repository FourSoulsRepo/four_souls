package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Magic Skin (Active Treasure Card)
//
//	{Tap Effect}Gain +1 treasure, put a counter on this, then give this to another player.
//	At the start of your turn, roll-
//	If the roll is less than the number of counters on this, destroy it and all other items you control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var magicSkin = engine.CardDef{
	Ref:    "magic_skin",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
