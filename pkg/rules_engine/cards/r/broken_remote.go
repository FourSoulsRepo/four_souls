package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Broken Remote (Trinket Card)
//
//	Each time you activate an item, you may put a room or monster not being attacked into discard.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var brokenRemote = engine.CardDef{
	Ref:    "broken_remote",
	Kind:   engine.LootCard,
	Copies: 1,
}
