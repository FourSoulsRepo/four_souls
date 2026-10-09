package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cain’s Eye (Trinket Card)
//
//	At the start of your turn, look at the top card of the loot deck. You may put it on the bottom.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cainsEye = engine.CardDef{
	Ref:    "cains_eye",
	Kind:   engine.LootCard,
	Copies: 1,
}
