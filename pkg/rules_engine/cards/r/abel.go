package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Abel (Passive Treasure Card)
//
//	Each time you deal combat damage to a monster, roll-
//	4-6: Deal 1 damage to another monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var abel = engine.CardDef{
	Ref:       "abel",
	Kind:      engine.TreasureCard,
	Copies:    1,
	ThreePlus: true,
}
