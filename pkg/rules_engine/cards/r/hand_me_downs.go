package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hand Me Downs (Active Treasure Card)
//
//	{Tap Effect}Each other player may choose to loot 1. Each player that does gives you a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var handMeDowns = engine.CardDef{
	Ref:    "hand_me_downs",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
