package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Book Of Virtues (Eternal Treasure Card)
//
//	Each time you deal combat damage, roll-
//	5-6: Deal 1 damage to a monster or player.
//	Each time you kill a monster, choose another player. You and that player each loot 1.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bookOfVirtues = engine.CardDef{
	Ref:     "book_of_virtues",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
