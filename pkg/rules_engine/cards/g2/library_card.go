package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Library Card (Active Treasure Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var libraryCard = engine.CardDef{
	Ref:    "library_card",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
