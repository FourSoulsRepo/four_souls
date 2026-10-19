package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Flip (Eternal Treasure Card)
//
//	Each time you die, after paying penalties, flip your character if able. Then recharge it. Discard your hand and loot 3.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var flip = engine.CardDef{
	Ref:     "flip",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
