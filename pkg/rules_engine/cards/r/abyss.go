package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Abyss (Eternal Treasure Card)
//
//	{Paid Effect}Destroy an item you control. If you do, put a counter on this.
//	You have +1{ATK} for every 2 counters on this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var abyss = engine.CardDef{
	Ref:     "abyss",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
