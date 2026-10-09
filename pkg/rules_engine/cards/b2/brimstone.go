package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Brimstone (Passive Treasure Card)
//
//	{ATK}
//	Each time you deal combat damage to a monster, deal 1 damage to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var brimstone = engine.CardDef{
	Ref:    "brimstone",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
