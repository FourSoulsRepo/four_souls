package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Belly Button (Passive Treasure Card)
//
//	You may play an additional loot card on your turn.
//	Each time you take damage, you may recharge your character.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bellyButton = engine.CardDef{
	Ref:    "belly_button",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
