package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blank Card (Active Treasure Card)
//
//	{Tap Effect}The next time you play a non-trinket, non-ambush loot card this turn, copy it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blankCard = engine.CardDef{
	Ref:    "blank_card",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
