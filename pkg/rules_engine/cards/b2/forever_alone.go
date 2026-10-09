package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Forever Alone (Eternal Treasure Card)
//
//	{Tap Effect}Choose one- Steal 1¢ from another player. Look at the top card of a deck. Discard a loot card, then loot 1.
//	Each time you take damage, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var foreverAlone = engine.CardDef{
	Ref:     "forever_alone",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
