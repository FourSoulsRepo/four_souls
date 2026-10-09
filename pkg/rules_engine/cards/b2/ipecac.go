package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ipecac (Passive Treasure Card)
//
//	{ATK}
//	Each time you roll an attack roll of 6, deal 1 damage to each other player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ipecac = engine.CardDef{
	Ref:    "ipecac",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
