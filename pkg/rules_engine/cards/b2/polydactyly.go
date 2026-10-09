package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Polydactyly (Passive Treasure Card)
//
//	You may play an additional loot card on your turn.
//	You have +1{ATK} for your first attack roll each turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var polydactyly = engine.CardDef{
	Ref:    "polydactyly",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
