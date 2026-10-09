package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Shadow (Passive Treasure Card)
//
//	If another player would pay the death penalty, you choose what item they would destroy and you gain any loot cards and ¢ they would lose.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var shadow = engine.CardDef{
	Ref:    "shadow",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
