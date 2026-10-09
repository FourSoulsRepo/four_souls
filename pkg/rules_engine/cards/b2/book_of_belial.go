package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Book Of Belial (Eternal Treasure Card)
//
//	Add or subtract 1 from a roll.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bookOfBelial = engine.CardDef{
	Ref:     "book_of_belial",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
