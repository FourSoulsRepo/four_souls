package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Berserk (Eternal Treasure Card)
//
//	You may attack an additional time on your turn.
//	Each time you kill a monster, gain +1{HP} till end of turn.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var berserk = engine.CardDef{
	Ref:     "berserk",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
