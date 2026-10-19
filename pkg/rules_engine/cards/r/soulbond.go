package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soulbond (Eternal Treasure Card)
//
//	Each time you die, choose another player. That player dies.
//	You don't lose ¢ or discard loot cards when paying the death penalty.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulbond = engine.CardDef{
	Ref:       "soulbond",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	ThreePlus: true,
	Outside:   true,
}
