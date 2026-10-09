package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Tarot Cloth (Passive Treasure Card)
//
//	Each time a player rolls a ❹, they must give you a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var tarotCloth = engine.CardDef{
	Ref:    "tarot_cloth",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
