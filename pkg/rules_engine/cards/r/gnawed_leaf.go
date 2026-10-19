package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Gnawed Leaf (Passive Treasure Card)
//
//	Prevent all damage you would take while it's not your turn.
//	The first time each turn another player plays a loot card that targets you or something you control, you may cancel it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gnawedLeaf = engine.CardDef{
	Ref:       "gnawed_leaf",
	Kind:      engine.TreasureCard,
	Copies:    1,
	ThreePlus: true,
}
