package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Glitch (Eternal Treasure Card)
//
//	As your turn starts, choose a shop item. This gains the abilities of that item till end of turn.
//	As your turn ends, choose a shop item. This gains the abilities of that item till the start of your next turn. Recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var glitch = engine.CardDef{
	Ref:     "glitch",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
