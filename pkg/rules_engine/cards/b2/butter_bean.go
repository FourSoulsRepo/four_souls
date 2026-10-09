package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Butter Bean! (Butter Bean Card)
//
//	Cancel the ↷ or $ ability of an item or a loot being played.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var butterBean = engine.CardDef{
	Ref:    "butter_bean",
	Kind:   engine.LootCard,
	Copies: 3,
}
