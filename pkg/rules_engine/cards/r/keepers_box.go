package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper’s Box (Passive Treasure Card)
//
//	When this enters play, expand shop slots by 2.
//	Each time another player purchases a shop item, gain 2¢ and loot 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var keepersBox = engine.CardDef{
	Ref:    "keepers_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
