package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Champion Belt (Passive Treasure Card)
//
//	You have +1{ATK} for your first attack roll each turn.
//	You may attack an additional time on your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var championBelt = engine.CardDef{
	Ref:    "champion_belt",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
