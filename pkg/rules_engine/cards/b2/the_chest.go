package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Chest (Soul Treasure Card)
//
//	if this would be destroyed, it becomes a soul instead.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theChest = engine.CardDef{
	Ref:    "the_chest",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
}
