package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lost Soul (Lost Soul Card)
//
//	When this enters play, it becomes a soul.
//	(It's no longer an item.)
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lostSoul = engine.CardDef{
	Ref:    "lost_soul",
	Kind:   engine.LootCard,
	Copies: 1,
	Soul:   1,
}
