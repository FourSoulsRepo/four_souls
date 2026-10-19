package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Tick (Trinket Card)
//
//	The first time you kill a monster on your turn, heal +1{HP} and you may attack an additional time this turn.
//	-Eternal- This can't be destroyed or put into discard.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var tick = engine.CardDef{
	Ref:     "tick",
	Kind:    engine.LootCard,
	Copies:  1,
	Eternal: true,
}
