package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pandora’s Box (Soul Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, roll-
//	1: Gain 1¢.
//	2: Gain 6¢.
//	3: Kill a monster.
//	4: Loot 3.
//	5: Gain 9¢.
//	6: This becomes a soul. Gain it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pandorasBox = engine.CardDef{
	Ref:    "pandoras_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
	Tap:    true,
}
