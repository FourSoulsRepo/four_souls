package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Yum Heart (Eternal Treasure Card)
//
//	{Tap Effect}Choose a player or monster. Prevent the next instance of damage they would take this turn.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var yumHeart = engine.CardDef{
	Ref:     "yum_heart",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
